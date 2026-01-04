import { FolderGit2, Download } from 'lucide-react';
import { Link } from 'react-router-dom';
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
        <Link to={`/resource/${row.original.user.name}/${row.original.path}`}
          className="flex flex-row justify-start items-center">
          <FolderGit2 className="px-1" />
          <div className="px-1 align-middle">{row.original.path.replace(/\.git/, '')}</div>
        </Link>
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
