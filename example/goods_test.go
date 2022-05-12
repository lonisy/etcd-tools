package example

import (
    "fmt"
    "testing"
)

func TestInitGoodsItems(t *testing.T) {
    InitGoodsItems()
    fmt.Println(len(GoodsItems))
    if len(GoodsItems) > 0 {
        fmt.Println(GoodsItems[0].GoodsID)
        fmt.Println(GoodsItems[0].Name)
        fmt.Println(GoodsItems[0].Image)
    }
}
