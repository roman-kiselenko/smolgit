import { useParams, Link } from 'react-router-dom';
import { useEffect } from 'react';
import { useReposState, getRepos } from '@/store/repositories';
import { toast } from 'sonner';
import { Repo } from '@/types';

export function RepoViewPage() {
  const { repoPath } = useParams<{ repoPath: string }>();
  const repos = useReposState();

  useEffect(() => {
    if (!repoPath) return;
    // fetch all repos if not already loaded
    if (!repos.repos.get().length) {
      getRepos('');
    }
  }, [repoPath]);

  const repoItem: Repo | any = repos.repos.get().find((r: Repo | any) => r.path === repoPath);

  if (!repoItem) {
    return (
      <div className="p-4">
        <h2 className="text-xl font-bold">Repository not found</h2>
        <Link to="/resource" className="text-blue-500">
          Back to list
        </Link>
      </div>
    );
  }

  return (
    <div className="flex-grow overflow-auto">
      <div className="flex flex-row py-2 px-2 items-center justify-between"></div>

      <div className="grid grid-cols-1">
        <div className="p-4">
          <h2 className="text-2xl font-bold mb-4">{repoItem.name}</h2>
          <p className="mb-2">
            <strong>Path:</strong> {repoItem.path}
          </p>
          <p className="mb-2">
            <strong>User:</strong> {repoItem.user?.name}
          </p>
          <Link to="/resource" className="text-blue-500">
            Back to list
          </Link>
        </div>
      </div>
    </div>
  );
}
