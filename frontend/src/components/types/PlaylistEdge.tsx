import { 
  BaseEdge, 
  type EdgeProps, 
  getBezierPath, 
  EdgeLabelRenderer 
} from '@xyflow/react';

export function PlaylistEdge({
  id,
  sourceX,
  sourceY,
  targetX,
  targetY,
  sourcePosition,
  targetPosition,
  style = {},
  markerEnd,
}: EdgeProps) {
  // 1. Generate the curved path
  const [edgePath, labelX, labelY] = getBezierPath({
    sourceX,
    sourceY,
    sourcePosition,
    targetX,
    targetY,
    targetPosition,
  });

  return (
    <>
      <path
        id={id + "_bg"}
        style={{
          fill: 'none',
          stroke: '#1DB954',
          strokeWidth: 6,
          opacity: 0.1, // Subtle green glow
        }}
        className="react-flow__edge-path"
        d={edgePath}
      />

      <BaseEdge
        path={edgePath}
        markerEnd={markerEnd}
        style={{
          ...style,
          stroke: '#1DB954', // Spotify Green
          strokeWidth: 2,
          transition: 'stroke 0.3s, stroke-width 0.3s',
        }}
      />

      {/* Optional: Add a small delete button in the middle of the edge */}
      <EdgeLabelRenderer>
        <div
          style={{
            position: 'absolute',
            transform: `translate(-50%, -50%) translate(${labelX}px,${labelY}px)`,
            pointerEvents: 'all',
          }}
          className="nodrag nopan"
        >
          {/* You can put an 'X' button here later to delete the inclusion */}
        </div>
      </EdgeLabelRenderer>
    </>
  );
}
