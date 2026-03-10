import { useCallback } from 'react';
import { 
  applyEdgeChanges, 
  applyNodeChanges, 
  ReactFlow, 
  Background, 
  Controls,
  type Edge, 
  type Node, 
  type EdgeChange, 
  type NodeChange,
  type Connection,
  MarkerType
} from '@xyflow/react';
import '@xyflow/react/dist/style.css';
import { usePlaylistStore } from '../stores/playlist.store';

import  CustomNode  from '../types/CustomNode'
import CustomConnectionLine from '../types/CustomConnectionLine';
import FloatingEdge from '../types/FloatingEdge';
import { createFlowEdge } from '../utils/flow-mapper';

export default function Flow() {
    const { 
        playlistNodeData, 
        playlistEdgeData, 
        setPlaylistNodes, 
        setPlaylistEdges,
        includePlaylist,
    } = usePlaylistStore();

    const nodeTypes = {
        custom: CustomNode,
    };

    const edgeTypes = {
        floating: FloatingEdge,
    };

    const connectionLineStyle = {
        stroke: '#b1b1b7',
    };

    const defaultEdgeOptions = {
        type: 'floating',
        markerEnd: {
            type: MarkerType.ArrowClosed,
            color: '#b1b1b7',
        },
    };

    const onNodesChange = useCallback(
        (changes: NodeChange<Node>[]) => {
            setPlaylistNodes((nds: Node[]) => applyNodeChanges(changes, nds));
        },
        [setPlaylistNodes],
    );

    const onEdgesChange = useCallback(
        (changes: EdgeChange<Edge>[]) => {
            setPlaylistEdges((eds: Edge[]) => applyEdgeChanges(changes, eds));
        },
        [setPlaylistEdges],
    );

    const onConnect = useCallback(
        (params: Connection) => {
            includePlaylist(params.target, params.source);
        },
        [setPlaylistEdges],
    );

    return (
        <div style={{ width: '100vw', height: '100vw' }} className="bg-[#121212]">
            <ReactFlow
                nodes={playlistNodeData}
                edges={playlistEdgeData}
                onNodesChange={onNodesChange} 
                onEdgesChange={onEdgesChange}
                onConnect={onConnect}
                colorMode="dark" 
                fitView
                nodeTypes={nodeTypes}
                edgeTypes={edgeTypes}
                defaultEdgeOptions={defaultEdgeOptions}
                connectionLineComponent={CustomConnectionLine}
                connectionLineStyle={connectionLineStyle}
            >
                <Background color="#282828" gap={20} />
                <Controls />
            </ReactFlow>
        </div>
    );
}
