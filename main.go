package main

import (
	"RedisExp/pkg"
	"fmt"
)

func main() {
	logo := `
██████╗ ███████╗██████╗ ██╗███████╗    ███████╗██╗  ██╗██████╗ 
██╔══██╗██╔════╝██╔══██╗██║██╔════╝    ██╔════╝╚██╗██╔╝██╔══██╗
██████╔╝█████╗  ██║  ██║██║███████╗    █████╗   ╚███╔╝ ██████╔╝
██╔══██╗██╔══╝  ██║  ██║██║╚════██║    ██╔══╝   ██╔██╗ ██╔═══╝ 
██║  ██║███████╗██████╔╝██║███████║    ███████╗██╔╝ ██╗██║     
╚═╝  ╚═╝╚══════╝╚═════╝ ╚═╝╚══════╝    ╚══════╝╚═╝  ╚═╝╚═╝
`
	fmt.Println(logo)

	pkg.Execute()

}
