package model

type Alert struct {
    ID          int64  `gorm:"primaryKey;column:id" json:"id"`
    TimeMs      int64  `gorm:"column:time_ms" json:"timeMs"`
    Type        string `gorm:"column:type" json:"type"` // offline, online, due
    NodeID      *int64 `gorm:"column:node_id" json:"nodeId,omitempty"`
    NodeName    *string `gorm:"column:node_name" json:"nodeName,omitempty"`
    Message     string `gorm:"column:message" json:"message"`
}

func (Alert) TableName() string { return "alert" }

