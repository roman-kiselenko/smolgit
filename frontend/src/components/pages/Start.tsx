import { useEffect, useCallback, useState } from 'react';
import { LogOut } from 'lucide-react';
import { DataTable } from '@/components/ui/DataTable';
import columns from '@/components/pages/Repos/Table/ColumnDef';
import { useReposState, getRepos } from '@/store/repositories';
import { Input } from '@/components/ui/input';
import { call } from '@/lib/api';
import { toast } from 'sonner';
import { Button } from '@/components/ui/button';
import { useAuth } from '@/context/AuthProvider';

export function StartPage() {
  const repos = useReposState();
  const [searchQuery, setSearchQuery] = useState('');
  const { logout, AuthDisabled } = useAuth();

  const fetchData = useCallback(async () => {
    try {
      await call<any[]>('ping');
      await getRepos(searchQuery);
    } catch (error: any) {
      toast.error('Error! Cant ping server\n' + error.message);
    }
  }, [searchQuery]);

  useEffect(() => {
    fetchData();

    const interval = setInterval(() => {
      fetchData();
    }, 2000);

    return () => clearInterval(interval);
  }, [fetchData]);

  return (
    <div className="flex-grow overflow-auto">
      <div className="flex flex-row py-2 px-2 items-center justify-between">
        <Input
          placeholder="Filter by name..."
          className="placeholder:text-muted-foreground flex h-6 w-full rounded-md bg-transparent py-2 text-sm outline-hidden disabled:cursor-not-allowed disabled:opacity-50"
          onChange={(e) => setSearchQuery(e.target.value)}
        />
        {!AuthDisabled && (
          <Button onClick={logout} className="ml-2 text-xs">
            <LogOut size={12} />
          </Button>
        )}
      </div>

      <div className="grid grid-cols-1">
        <div className="mx-3 h-24 col-span-2">
          <DataTable
            menuDisabled={true}
            kind={'repos'}
            noResult={true}
            columns={columns as any}
            data={repos.repos.get() as any}
          />
        </div>
      </div>
    </div>
  );
}
