import { Handle, Position } from '@xyflow/react';

export function PlaylistNode({ data }: { data: any }) {
  return (
    <div className="bg-[#181818] text-white border border-[#282828] rounded-lg p-3 text-[12px] w-[180px] text-center shadow-xl group hover:border-[#1DB954] transition-colors">
      <Handle 
        type="source" 
        position={Position.Top} 
        className="w-3 h-3 bg-[#1DB954] border-2 border-[#121212] !-top-1.5" 
      />
      
      <div className="font-bold truncate">{data.label}</div>
      
      <Handle 
        type="target" 
        position={Position.Bottom} 
        className="w-3 h-3 bg-gray-500 border-2 border-[#121212] !-bottom-1.5" 
      />
    </div>
  );
}
