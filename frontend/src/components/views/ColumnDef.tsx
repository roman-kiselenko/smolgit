import { Clock } from 'lucide-react';
import AgeCell from '@/components/ui/Table/AgeCell';
import { ColumnDef } from '@tanstack/react-table';

const columns: ColumnDef<string>[] = [
  {
    accessorKey: '',
    id: 'name',
    header: 'File',
    cell: ({ row }) => {
      return <div className="px-1 align-middle">{row.original.filename}</div>;
    },
  },
  {
    accessorKey: 'time',
    id: 'time',
    header: ({ column }) => <Clock size={12} />,
    cell: ({ row }) => {
      return (
        <div className="flex flex-row justify-start">
          <AgeCell age={'Nov 9, 2025, 11:38 AM GMT+3'} />
        </div>
      );
    },
  },
];

export default columns;
