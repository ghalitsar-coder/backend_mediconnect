package main
import (
	"fmt"
	"gorm.io/gorm/schema"
)
func main() {
    ns := schema.NamingStrategy{}
	fmt.Printf("Col: %s\n", ns.ColumnName("", "KtpURL"))
}
