import {
  ColumnDef,
  flexRender,
  getCoreRowModel,
  SortingState,
  RowSelectionState,
  getFilteredRowModel,
  getSortedRowModel,
  useReactTable,
} from '@tanstack/react-table';
import { useState } from 'react';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';

import { useNavigate } from 'react-router';

declare module '@tanstack/react-table' {
  /* eslint-disable @typescript-eslint/no-unnecessary-type-constraint */
  /* eslint-disable @typescript-eslint/no-unused-vars */
  interface ColumnMeta<TData extends unknown, TValue> {
    className?: string;
  }
}

interface ColumnData {
  kind: string;
  apiVersion: string;
  metadata?: any; // or a more specific type for metadata
}

interface DataTableProps<TData, TValue> {
  kind: string;
  columns: ColumnDef<TData, TValue>[];
  data: TData[];
  noResult?: boolean;
  doubleClickDisabled?: boolean;
  menuDisabled?: boolean;
}

export function DataTable<TData, TValue>({
  kind,
  columns,
  data,
  noResult,
  doubleClickDisabled = false,
  menuDisabled = false,
}: DataTableProps<TData, TValue>) {
  const [sorting, setSorting] = useState<SortingState>([]);
  const [rowSelection, setRowSelection] = useState<RowSelectionState>({});

  const table = useReactTable({
    data,
    columns,
    getCoreRowModel: getCoreRowModel(),
    onSortingChange: setSorting,
    onRowSelectionChange: setRowSelection,
    getSortedRowModel: getSortedRowModel(),
    getFilteredRowModel: getFilteredRowModel(),
    getRowId: (row) => (row as ColumnData).metadata?.uid,
    state: {
      sorting,
      rowSelection,
    },
  });
  const tableRows = table
    .getSelectedRowModel()
    .rows.filter((x) => (x.original as ColumnData)?.metadata);
  return (
    <Table>
      <TableHeader className="text-xs">
        {table.getHeaderGroups().map((headerGroup) => (
          <TableRow key={headerGroup.id}>
            {headerGroup.headers.map((header) => {
              return (
                <TableHead key={header.id}>
                  {header.isPlaceholder
                    ? null
                    : flexRender(header.column.columnDef.header, header.getContext())}
                </TableHead>
              );
            })}
          </TableRow>
        ))}
      </TableHeader>
      <TableBody className="text-xs">
        {table.getRowModel().rows?.length ? (
          table.getRowModel().rows.map((row) => (
            <TableRow key={row.id} data-state={row.getIsSelected() && 'selected'}>
              {row.getVisibleCells().map((cell, i) => (
                <TableCell
                  className={cell.column.columnDef.meta?.className as string}
                  key={i}
                  onDoubleClick={() => console.log('click')}
                  onClick={row.getToggleSelectedHandler()}
                >
                  {flexRender(cell.column.columnDef.cell, cell.getContext())}
                </TableCell>
              ))}
            </TableRow>
          ))
        ) : (
          <TableRow>
            <TableCell colSpan={columns.length} className="h-24 text-center">
              {noResult ? 'No results' : <></>}
            </TableCell>
          </TableRow>
        )}
      </TableBody>
    </Table>
  );
}
