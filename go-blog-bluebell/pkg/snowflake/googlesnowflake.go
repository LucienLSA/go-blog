package snowflake

import (
	"time"

	"github.com/bwmarrin/snowflake"
)

var node *snowflake.Node

func InitSnowflake(startTime string, machineID int64) (err error) {
	var st time.Time
	st, err = time.Parse("2006-01-02", startTime)
	if err != nil {
		return
	}
	snowflake.Epoch = st.UnixNano() / 1000000
	node, err = snowflake.NewNode(machineID)
	return
}

func GenID() int64 {
	return node.Generate().Int64()
}

// func main() {
// 	if err := InitSnowflake("2020-07-01", 1); err != nil {
// 		fmt.Printf("init failed, err:%v\n", err)
// 		return
// 	}
// 	id := GenID()
// 	fmt.Println(id)
// }
