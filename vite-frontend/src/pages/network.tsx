import { useEffect, useMemo, useState } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { Button } from "@heroui/button";
import { Card, CardBody, CardHeader } from "@heroui/card";
import { getNodeNetworkStats, getNodeNetworkStatsBatch, getNodeList } from "@/api";
import { useRef } from "react";
import toast from "react-hot-toast";

const ranges = [
  { key: '1h', label: '每小时' },
  { key: '12h', label: '每12小时' },
  { key: '1d', label: '每天' },
  { key: '7d', label: '每七天' },
  { key: '30d', label: '每月' },
];

export default function NetworkPage() {
  const params = useParams();
  const navigate = useNavigate();
  const nodeId = Number(params.id);
  const [range, setRange] = useState('1h');
  const [data, setData] = useState<any>({ results: [], targets: {}, disconnects: [], sla: 0 });
  const [nodes, setNodes] = useState<any[]>([]);
  const [batch, setBatch] = useState<any>({});
  const [nodeName, setNodeName] = useState<string>("");
  const [loading, setLoading] = useState(false);
  const chartRef = useRef<HTMLDivElement>(null);
  const chartInstanceRef = useRef<any>(null);

  const load = async () => {
    setLoading(true);
    try {
      if (params.id) {
        const res = await getNodeNetworkStats(nodeId, range);
        if (res.code === 0) setData(res.data || { results: [], disconnects: [], sla: 0 });
        else toast.error(res.msg || '加载失败');
      } else {
        const [l, b] = await Promise.all([getNodeList(), getNodeNetworkStatsBatch(range)]);
        if (l.code === 0) setNodes(l.data || []);
        if (b.code === 0) setBatch(b.data || {});
      }
    } catch { toast.error('网络错误'); } finally { setLoading(false); }
  };
  useEffect(() => { load(); }, [params.id, range]);

  // fetch node name for detail page
  useEffect(() => {
    if (params.id) {
      getNodeList().then((res:any)=>{
        if (res.code===0 && Array.isArray(res.data)){
          const n = res.data.find((x:any)=>x.id===Number(params.id));
          if (n) setNodeName(n.name||`节点 ${params.id}`);
        }
      }).catch(()=>{});
    } else {
      setNodeName("");
    }
  }, [params.id]);

  const grouped = useMemo(() => {
    const g: Record<string, any[]> = {};
    for (const r of (data.results || [])) {
      const k = String(r.targetId);
      (g[k] ||= []).push(r);
    }
    return g;
  }, [data]);

  useEffect(() => {
    const render = async () => {
      if (!chartRef.current) return;
      const echarts = await import('echarts');
      if (!chartInstanceRef.current) {
        chartInstanceRef.current = echarts.init(chartRef.current);
      }
      const series: any[] = [];
      Object.keys(grouped).forEach((tid) => {
        const arr = grouped[tid];
        const label = data.targets?.[tid]?.name || `目标${tid}`;
        series.push({
          type: 'line', sampling: 'lttb',
          name: `${label} RTT`,
          showSymbol: false,
          yAxisIndex: 0,
          data: arr.map((it:any)=>[it.timeMs, it.ok? it.rttMs : null])
        });
        series.push({
          type: 'line', sampling: 'lttb',
          name: `${label} 丢包%`,
          showSymbol: false,
          yAxisIndex: 1,
          data: arr.map((it:any)=>[it.timeMs, it.ok? 0 : 100])
        });
      });
      chartInstanceRef.current.setOption({
        tooltip: { trigger: 'axis' },
        legend: { type: 'scroll' },
        dataZoom: [
          { type: 'inside', throttle: 50 },
          { type: 'slider', height: 20 }
        ],
        xAxis: { type: 'time' },
        yAxis: [
          { type: 'value', name: 'RTT (ms)' },
          { type: 'value', name: '丢包(%)', min: 0, max: 100, axisLabel: { formatter: '{value}%' } }
        ],
        series,
        grid: { left: 40, right: 20, top: 40, bottom: 30 }
      });
      window.addEventListener('resize', handleResize);
    };
    const handleResize = () => { try { chartInstanceRef.current?.resize(); } catch {} };
    render();
    return () => { window.removeEventListener('resize', handleResize); };
  }, [grouped, data.targets]);

  return (
    <div className="px-4 py-6 space-y-4">
      <div className="flex items-center justify-between">
        {params.id ? (
          <>
            <h2 className="text-xl font-semibold">{nodeName || `节点 ${params.id}`} 网络详情</h2>
            <div className="text-sm text-default-500">SLA：{(data.sla*100).toFixed(2)}%</div>
          </>
        ) : (
          <h2 className="text-xl font-semibold">节点网络概览</h2>
        )}
      </div>

      <div className="flex gap-2">
        {ranges.map(r => (
          <Button key={r.key} size="sm" variant={range===r.key? 'solid':'flat'} color={range===r.key? 'primary':'default'} onPress={()=>setRange(r.key)}>
            {r.label}
          </Button>
        ))}
      </div>

      {params.id ? (
      <Card>
        <CardHeader className="justify-between">
          <div className="font-semibold">Ping 统计（按目标）</div>
          <Button size="sm" variant="flat" onPress={load} isLoading={loading}>刷新</Button>
        </CardHeader>
        <CardBody>
          <div className="h-[360px]" ref={chartRef} />
        </CardBody>
      </Card>
      ) : (
      <Card>
        <CardHeader className="justify-between">
          <div className="font-semibold">节点网络概览（{range}）</div>
          <Button size="sm" variant="flat" onPress={load} isLoading={loading}>刷新</Button>
        </CardHeader>
        <CardBody>
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3">
            {nodes.map((n:any)=>{
              const s = batch?.[n.id] || {};
              const avg = s.avg ?? null; const latest = s.latest ?? null;
              return (
                <div key={n.id} className="p-3 rounded border border-default-200">
                  <div className="font-semibold mb-1">{n.name}</div>
                  <div className="text-sm text-default-500">最新: {latest!=null? `${latest} ms` : '-'}</div>
                  <div className="text-sm text-default-500">平均: {avg!=null? `${avg} ms` : '-'}</div>
                  <div className="mt-2"><Button size="sm" variant="flat" onPress={()=>navigate(`/network/${n.id}`)}>查看详情</Button></div>
                </div>
              )
            })}
          </div>
        </CardBody>
      </Card>
      )}

      {params.id && (
      <Card>
        <CardHeader className="font-semibold">断联记录</CardHeader>
        <CardBody>
          <div className="space-y-2 text-sm">
            {(data.disconnects || []).map((it:any)=>{
              const dur = it.durationS ? it.durationS : (it.upAtMs ? Math.round((it.upAtMs - it.downAtMs)/1000) : null);
              return (
                <div key={it.id} className="flex justify-between p-2 rounded bg-default-50">
                  <div>开始：{new Date(it.downAtMs).toLocaleString()}</div>
                  <div>恢复：{it.upAtMs ? new Date(it.upAtMs).toLocaleString() : '-'}</div>
                  <div>时长：{dur !== null ? `${dur}s` : '-'}</div>
                </div>
              )
            })}
            {(!data.disconnects || data.disconnects.length===0) && <div className="text-default-500">暂无记录</div>}
          </div>
        </CardBody>
      </Card>
      )}
    </div>
  );
}
