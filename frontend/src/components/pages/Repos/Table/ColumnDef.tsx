import { FolderGit2 } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { ColumnDef } from '@tanstack/react-table';
import { Repo } from '@/types';
const columns: ColumnDef<Repo>[] = [
  {
    accessorKey: 'repo',
    id: 'path',
    header: 'Repo',
    cell: ({ row }) => {
      return <div className="flex flex-row justify-start">
        <FolderGit2 className="px-1"/>
        <div className="px-1 align-middle">{row.original.path}</div>
      </div>;
    },
  },
  {
    accessorKey: 'user.name',
    id: 'user',
    header: 'User',
    cell: ({ row }) => {
      return <div>{row.original?.user.name}</div>;
    },
  },

];

export default columns;
