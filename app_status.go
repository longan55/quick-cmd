package main

// Status 保存应用初始状态
type Status struct {
	Os        []OS      `json:"os"`
	SortIndex SortIndex `json:"sortIndex"`
}

// SortIndex 保存排序配置
type SortIndex struct {
	CreatTimeAsc bool `json:"creatTimeAsc"`
	IDAsc        bool `json:"idAsc"`
	NameAsc      bool `json:"nameAsc"`
	SortValueAsc bool `json:"sortValueAsc"`
	CopyCountAsc bool `json:"copyCountAsc"`
}

func (a *App) GetStatus() *Status {
	return &Status{
		Os: []OS{
			Linux,
		},
		SortIndex: SortIndex{
			CreatTimeAsc: true,
			IDAsc:        true,
			NameAsc:      true,
			SortValueAsc: true,
			CopyCountAsc: true,
		},
	}
}
