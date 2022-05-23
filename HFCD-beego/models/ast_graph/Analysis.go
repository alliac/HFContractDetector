package ast_graph

import (
	"HFCD-beego/models/ast_graph/bugs"
	"HFCD-beego/models/ast_graph/gen"
	"time"
)

func Main_1(ReportPath string) time.Duration {
	start := time.Now()
	bugs.WriteToReport(ReportPath, "---------------------------Static Analysis---------------------------")
	path := "./models/chaincode.go"
	dpath := "./models/ast_graph"
	bugs.InitDtct(ReportPath)
	gen.GenSvg(path, dpath, "tree")
	bugs.FinshDtct(path, ReportPath)
	//bugs.WriteToReport(ReportPath, "---------------------------Static Analysis finished!---------------------------")
	var cost time.Duration
	cost = time.Since(start)
	//bugs.WriteToReport(ReportPath, "cost="+fmt.Sprint(cost))
	return cost
}
