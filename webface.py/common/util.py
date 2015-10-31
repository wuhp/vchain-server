class VisNode(object):
    def __init__(self):
        self.name = None
        self.x = 0
        self.y = 0


class VisEdge(object):
    def __init__(self):
        self.source = None
        self.target = None
        self.name = None
        self.path = None


class VisGraph(object):
    def __init__(self):
        self.width = 0
        self.height = 0
        self.nodes = []
        self.edges = []


def graphNormalize(graph, w, h):
    if w <= 0 or h <=0:
        raise ValueError("width or height should not be 0")
    n = len(graph.nodes)
    if n == 0:
        return
    if n == 1:
        node = graph.nodes[0]
        node.x = w / 2
        node.y = h / 2
        return
    count = [0] * n
    cindex = [0] * n
    yslice = [0] * n
    level = 0
    for node in graph.nodes:
        if level <= node.x:
            level = node.x + 1
        count[node.x] += 1
    xslice = w / 2.0 / level
    for i, one in enumerate(count):
        if one == 0:
            continue
        yslice[i] = h / 2.0 / one
    for node in graph.nodes:
        level = node.x
        node.y = int(yslice[level]+yslice[level]*2*cindex[level])
        node.x = int(xslice+xslice*2*level)
        cindex[level] += 1


def graphBuildInvokeChain(graph, names, parents):
    n = len(names)
    if n != len(parents):
        raise ValueError("the length of name array and parent array do not match")
    if parents[0] != -1:
        raise ValueError("the first node is not root")
    graph.width = 0
    graph.height = 0
    graph.edges = []
    graph.nodes = []
    for i, one in enumerate(names):
        node = VisNode()
        node.name = one
        graph.nodes.append(node)
        if i <= 0: # skip first node
            continue
        edge = VisEdge()
        edge.source = graph.nodes[parents[i]]
        edge.target = node
        edge.name = str(i)
        graph.edges.append(edge)
        # currently X store as level
        node.x = edge.source.x + 1


def svgGenerate(graph):
    result = [
        """<svg width="400" height="200" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink">""",
        """ <defs>""",
        """  <marker id="markerArrow" markerWidth="13" markerHeight="13" refX="2" refY="6" orient="auto" markerUnits="userSpaceOnUse">""",
        """   <path d="M2,2 L2,11 L10,6 L2,2" style="fill: black;" />""",
        """  </marker>""",
        """ </defs>""",
    ]
    for i, edge in enumerate(graph.edges):
        edge.path = "M%d,%d L%d,%d L%d,%d" % (edge.source.x, edge.source.y, (edge.source.x + edge.target.x)/2, (edge.source.y + edge.target.y)/2, edge.target.x, edge.target.y)
        result.append(""" <path id="edge-%d" d="%s" style="stroke: black; marker-mid: url(#markerArrow)" />""" % (i, edge.path))
    for i, node in enumerate(graph.nodes):
        result.append(""" <circle id="node-%d" cx="%d" cy="%d" r="15" style="fill:yellow; stroke:black;" />""" % (i, node.x, node.y))
        result.append(""" <text x="%d" y="%d" style="text-anchor:middle; font-size:12px; fill:black;">%s</text>""" % (node.x, node.y, node.name))
    result.append("</svg>")
    return "\n".join(result)


# e.g:
# names = ["A", "B", "C", "D", "B", "C", "E"]
# parents = [-1, 0, 1, 1, 3, 3, 4]
# graph = VisGraph()
# graphBuildInvokeChain(graph, names, parents)
# graphNormalize(graph, 400, 200)
# print svgGenerate(graph)
