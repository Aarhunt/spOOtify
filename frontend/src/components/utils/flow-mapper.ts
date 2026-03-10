import type { ModelPlaylistInclusionResponse, ModelPlaylistResponse } from '@/client';
import dagre from '@dagrejs/dagre';
import { type Node, type Edge, Position, MarkerType, type Connection } from '@xyflow/react';

const dagreGraph = new dagre.graphlib.Graph();
dagreGraph.setDefaultEdgeLabel(() => ({}));

const nodeWidth = 185;
const nodeHeight = 84;

export const getLayoutedElements = (nodes: Node[], edges: Edge[]) => {
  dagreGraph.setGraph({ rankdir: 'LR', nodesep: 40, ranksep: 100 });

  nodes.forEach((node) => {
    dagreGraph.setNode(node.id, { width: nodeWidth, height: nodeHeight });
  });

  edges.forEach((edge) => {
    dagreGraph.setEdge(edge.source, edge.target);
  });

  dagre.layout(dagreGraph);

  const layoutedNodes = nodes.map((node) => {
    const nodeWithPosition = dagreGraph.node(node.id);
    return {
      ...node,
      targetPosition: Position.Top,
      sourcePosition: Position.Bottom,
      position: {
        x: nodeWithPosition.x - nodeWidth / 2,
        y: nodeWithPosition.y - nodeHeight / 2,
      },
    };
  });

  return { nodes: layoutedNodes, edges };
};

export const mapPlaylistsToNodes = (playlists: ModelPlaylistResponse[]): Node[] => {
  return playlists.map((p) => ({
    id: p.spotifyID!,
    type: 'custom',
    data: { label: p.name },
    position: { x: 0, y: 0 }, 
  }));
};

export const mapInclusionsToEdges = (inclusions: ModelPlaylistInclusionResponse[]): Edge[] => {
    return inclusions.map((el) => createFlowEdge({ 
        source: el.source!, 
        target: el.target! 
    }));
};

export const createFlowEdge = (params: Connection | { source: string; target: string }): Edge => {
    return {
        ...params,
        id: `e-${params.source}-${params.target}`,
        type: 'floating',
        animated: true,
        markerEnd: {
            type: MarkerType.ArrowClosed,
            color: '#b1b1b7',
        },
    } as Edge;
};
