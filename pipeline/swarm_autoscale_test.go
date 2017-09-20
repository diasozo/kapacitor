package pipeline

import (
	"bytes"
	"testing"
	"time"

	"github.com/influxdata/kapacitor/tick/ast"
)

func TestSwarmAutoScaleNode_Tick(t *testing.T) {
	type args struct {
		cluster              string
		serviceName          string
		serviceNameTag       string
		outputServiceNameTag string
		currentField         string
		max                  int64
		min                  int64
		replicas             *ast.LambdaNode
		increaseCooldown     time.Duration
		decreaseCooldown     time.Duration
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "upgrade mutalisk to guardian",
			args: args{
				cluster:              "zerg",
				serviceName:          "units",
				serviceNameTag:       "mutalisk",
				outputServiceNameTag: "guardian",
				currentField:         "hitPoints",
				max:                  120,
				min:                  0,
				replicas: &ast.LambdaNode{
					Expression: &ast.FunctionNode{
						Type: ast.GlobalFunc,
						Func: "if",
						Args: []ast.Node{
							&ast.BinaryNode{
								Operator: ast.TokenGreater,
								Left: &ast.ReferenceNode{
									Reference: "greater spire",
								},
								Right: &ast.NumberNode{
									IsInt: true,
									Int64: 1,
									Base:  10,
								},
							},
							&ast.BinaryNode{
								Operator: ast.TokenPlus,
								Left: &ast.ReferenceNode{
									Reference: "replicas",
								},
								Right: &ast.NumberNode{
									IsInt: true,
									Int64: 1,
									Base:  10,
								},
							},
							&ast.ReferenceNode{
								Reference: "replicas",
							},
						},
					},
				},
				increaseCooldown: 6670 * time.Millisecond,
				decreaseCooldown: 2500 * time.Millisecond,
			},
			want: `|swarmAutoscale().cluster('zerg').serviceName('units').serviceNameTag('mutalisk').outputServiceNameTag('guardian').currentField('hitPoints').max(120).replicas(lambda: if("greater spire" > 1, "replicas" + 1, "replicas")).increaseCooldown(6670ms).decreaseCooldown(2500ms)`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := newSwarmAutoscaleNode(StreamEdge)
			n.Cluster = tt.args.cluster
			n.ServiceName = tt.args.serviceName
			n.ServiceNameTag = tt.args.serviceNameTag
			n.OutputServiceNameTag = tt.args.outputServiceNameTag
			n.CurrentField = tt.args.currentField
			n.Max = tt.args.max
			n.Min = tt.args.min
			n.Replicas = tt.args.replicas
			n.IncreaseCooldown = tt.args.increaseCooldown
			n.DecreaseCooldown = tt.args.decreaseCooldown

			var buf bytes.Buffer
			n.Tick(&buf)
			got := buf.String()
			if got != tt.want {
				t.Errorf("%q. TestSwarmAutoscaleNode_Tick() =\n%v\n want\n%v\n", tt.name, got, tt.want)
			}
		})
	}
}
