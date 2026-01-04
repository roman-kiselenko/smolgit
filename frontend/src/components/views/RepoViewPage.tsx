import { useParams, Link } from 'react-router-dom';
import { LogOut } from 'lucide-react';
import { useEffect, useState } from 'react';
import { useReposState, getRepos } from '@/store/repositories';
import { useRepoFilesState, getRepoFiles } from '@/store/repofiles';
import { Repo } from '@/types';
import { DataTable } from '@/components/ui/DataTable';
import columns from '@/components/views/ColumnDef';
import { Input } from '@/components/ui/input';
import { useAuth } from '@/context/AuthProvider';
import { Button } from '@/components/ui/button';

export function RepoViewPage() {
  const { user, repoPath } = useParams<{ user: string; repoPath: string }>();
  const repos = useReposState();
  const repoFiles = useRepoFilesState();
  const [searchQuery, setSearchQuery] = useState('');
  const { logout, AuthDisabled } = useAuth();

  useEffect(() => {
    if (!repoPath) return;
    if (!repos.repos.get().length) {
      getRepos('');
    }
  }, [repoPath]);

  useEffect(() => {
    if (!user || !repoPath) return;
    getRepoFiles(searchQuery, user, repoPath);
  }, [user, searchQuery]);

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
        <div className="mx-3">
          <DataTable
            menuDisabled={true}
            kind="repo"
            noResult={false}
            columns={columns as any}
            data={repoFiles.files.get() as any}
          />
        </div>
      </div>
    </div>
  );
}
