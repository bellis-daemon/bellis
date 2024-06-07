package assets

import _ "embed"

//go:embed echarts.min.js
var EChartsMinJs []byte

//go:embed echarts@4.min.js
var ECharts4MinJs []byte
