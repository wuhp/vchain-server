package webface;

import (
  "fmt"
  "strconv"
);

type VisualNode struct {
  Reference interface{};
  Title     string;
  X, Y      int;
}

type VisualEdge struct {
  Reference interface{};
  From, To  *VisualNode;
  Title     string;
  Path      string;
}

type VisualGraph struct {
  Edges  []*VisualEdge;
  Nodes  []*VisualNode;
  Width  int;
  Height int;
}

func (node *VisualNode) ToString () string {
  return fmt.Sprintf("%+v", node);
}

func (edge *VisualEdge) ToString () string {
  return fmt.Sprintf("%+v", edge);
}

func (graph *VisualGraph) Normalize (w int, h int) {
  if w <= 0 || h <= 0 {
    panic("width or height must be positive.");
  }
  n := len(graph.Nodes)

  if n == 0 { return; }
  if n == 1 {
    node := graph.Nodes[0]
    node.X = w/2;
    node.Y = h/2;
    return;
  }

  count := make([]int, n, n)
  cindex := make([]int, n, n)
  yslice := make([]float64, n, n)
  // TODO: fill 0 to each element in `count`
  level := 0
  for _, node := range graph.Nodes {
    if level <= node.X { level = node.X + 1; }
    count[node.X] ++;
  }

  xslice := float64(w)/float64(level*2);
  for i, one := range count {
    if one == 0 { continue; }
    yslice[i] = float64(h)/float64(2*one);
  }

  for _, node := range graph.Nodes {
    level = node.X;
    node.Y = int(yslice[level]+yslice[level]*float64(2*cindex[level]))
    node.X = int(xslice+xslice*float64(2*level));
    cindex[level] ++;
  }
}

func (graph *VisualGraph) BuildInvokeChain (name []string, parent []int) {
  n := len(name)
  if n != len(parent) {
    panic("the length of name array and parent array do not match.");
  }
  if parent[0] != -1 {
    panic("the first node is not root.");
  }
  graph.Width = 0;
  graph.Height = 0;
  graph.Edges = make([]*VisualEdge, n-1, n-1);
  graph.Nodes = make([]*VisualNode, n, n);

  for i, one := range name {
    node := new(VisualNode);
    node.X = 0;
    node.Y = 0;
    node.Title = one;
    graph.Nodes[i] = node;

    // skip first node
    if i <= 0 { continue; }
    edge := new(VisualEdge);
    // parent[i] must less than i
    edge.From = graph.Nodes[parent[i]];
    edge.To = node;
    edge.Title = strconv.Itoa(i)
    graph.Edges[i-1] = edge;

    // current X store as level
    node.X = edge.From.X + 1;
  }

}

func GenerateSVG(graph *VisualGraph) string {
  result := "";
  result += fmt.Sprintf(`<svg width="400" height="200" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink">`);
  result += fmt.Sprintf(` <defs>`);
  result += fmt.Sprintf(`  <marker id="markerArrow" markerWidth="13" markerHeight="13" refX="2" refY="6" orient="auto" markerUnits="userSpaceOnUse">`);
  result += fmt.Sprintf(`   <path d="M2,2 L2,11 L10,6 L2,2" style="fill: black;" />`);
  result += fmt.Sprintf(`  </marker>`);
  result += fmt.Sprintf(` </defs>`);
  for i, edge := range graph.Edges {
    edge.Path = fmt.Sprintf(
      "M%d,%d L%d,%d L%d,%d", edge.From.X, edge.From.Y, (edge.From.X + edge.To.X)/2, (edge.From.Y + edge.To.Y)/2, edge.To.X, edge.To.Y);
    result += fmt.Sprintf(` <path id="edge-%d" d="%s" style="stroke: black; marker-mid: url(#markerArrow)" />`, i, edge.Path);
  }
  for i, node := range graph.Nodes {
    result += fmt.Sprintf(` <circle id="node-%d" cx="%d" cy="%d" r="15" style="fill:yellow; stroke:black;" />`, i, node.X, node.Y);
    result += fmt.Sprintf(` <text x="%d" y="%d" style="text-anchor:middle; font-size:12px; fill:black;">%s</text>`, node.X, node.Y, node.Title);
  }
  result += fmt.Sprintf(`</svg>`);
  return result;
}

func test_dump (graph *VisualGraph) {
  fmt.Println(`<?xml version="1.0"?>`);
  fmt.Println(`<svg width="400" height="200" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink">`);
  fmt.Println(` <defs>`);
  fmt.Println(`  <marker id="markerArrow" markerWidth="13" markerHeight="13" refX="2" refY="6" orient="auto" markerUnits="userSpaceOnUse">`);
  fmt.Println(`   <path d="M2,2 L2,11 L10,6 L2,2" style="fill: black;" />`);
  fmt.Println(`  </marker>`);
  fmt.Println(` </defs>`);
  for _, edge := range graph.Edges {
    edge.Path = fmt.Sprintf(
      "M%d,%d L%d,%d L%d,%d", edge.From.X, edge.From.Y, (edge.From.X + edge.To.X)/2, (edge.From.Y + edge.To.Y)/2, edge.To.X, edge.To.Y);
    fmt.Printf(` <path d="%s" style="stroke: black; marker-mid: url(#markerArrow)" />`, edge.Path);
    fmt.Println();
  }
  for _, node := range graph.Nodes {
    fmt.Printf(` <circle cx="%d" cy="%d" r="15" style="fill:yellow; stroke:black;" />`, node.X, node.Y);
    fmt.Printf(` <text x="%d" y="%d" style="text-anchor:middle; font-size:12px; fill:black;">%s</text>`, node.X, node.Y, node.Title);
    fmt.Println();
  }
  fmt.Println(`</svg>`);
}

func TestMain () {
  /*
            C
            ^
           /
     A -> B   B -> E
           \  ^
            v/
            D
   */
  name := []string{"A", "B", "C", "D", "B", "C", "E"};
  parent := []int{-1, 0, 1, 1, 3, 3, 4};
  graph := new(VisualGraph);
  graph.BuildInvokeChain(name, parent);
  graph.Normalize(400, 200);
  test_dump(graph);
}
