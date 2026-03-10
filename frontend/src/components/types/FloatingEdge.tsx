import { BaseEdge, EdgeLabelRenderer, getStraightPath, useInternalNode, type EdgeProps } from '@xyflow/react';

import { getEdgeParams } from '../utils/utils';
import { usePlaylistStore } from '../stores/playlist.store';

function FloatingEdge({ id, source, target, markerEnd, style }: EdgeProps) {
  const sourceNode = useInternalNode(source);
  const targetNode = useInternalNode(target);

  if (!sourceNode || !targetNode) {
    return null;
  }

  const { setPlaylistEdges, undoIncludePlaylist } = usePlaylistStore();

  const { sx, sy, tx, ty } = getEdgeParams(sourceNode, targetNode);

  const [path, labelX, labelY] = getStraightPath({
    sourceX: sx,
    sourceY: sy,
    targetX: tx,
    targetY: ty,
  });

  const onEdgeClick = () => {
    setPlaylistEdges((edges) => edges.filter((edge) => edge.id !== id));
    undoIncludePlaylist(target, source)
  };

  return (
      <>
    <BaseEdge
      id={id}
      className="react-flow__edge-path"
      path={path}
      markerEnd={markerEnd}
      style={style}
    />
      <EdgeLabelRenderer>
        <div
          className="button-edge__label nodrag nopan"
          style={{
            transform: `translate(-50%, -50%) translate(${labelX}px,${labelY}px)`,
          }}
        >
          <button className="button-edge__button" onClick={onEdgeClick}>
            ×
          </button>
        </div>
      </EdgeLabelRenderer>
      </>
  );
}

export default FloatingEdge;
