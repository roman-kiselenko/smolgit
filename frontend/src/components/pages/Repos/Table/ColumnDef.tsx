import { FolderGit2, Download } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { ColumnDef } from '@tanstack/react-table';
import { Repo } from '@/types';
const columns: ColumnDef<Repo>[] = [
  {
    accessorKey: 'repo',
    id: 'path',
    header: 'Repo',
    cell: ({ row }) => {
      return (
        <div className="flex flex-row justify-start">
          <FolderGit2 className="px-1" />
          <div className="px-1 align-middle">{row.original.path.replace(/\.git/, '')}</div>
        </div>
      );
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
  {
    accessorKey: 'download',
    id: 'download',
    header: '',
    cell: ({ row }) => {
      return (
        <div className="flex flex-row justify-start">
          <Button>
            <Download size={12} />
          </Button>
        </div>
      );
    },
  },
];

export default columns;
