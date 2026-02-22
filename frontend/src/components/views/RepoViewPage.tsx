import { Button } from '@/components/ui/button';
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from '@/components/ui/collapsible';
import { ChevronRightIcon, FileIcon, FolderIcon } from 'lucide-react';
import { useParams, Link } from 'react-router-dom';
import { LogOut } from 'lucide-react';
import { useEffect, useState } from 'react';
import { useReposState, getRepos } from '@/store/repositories';
import { useRepoFilesState, getRepoFiles } from '@/store/repofiles';
import { Repo } from '@/types';
import { Input } from '@/components/ui/input';
import { useAuth } from '@/context/AuthProvider';
import { CardContent, CardHeader } from '@/components/ui/card';
import timeAgo from '@/timeAgo';
import moment from 'moment';

type FileTreeItem = { name: string } | { name: string; items: FileTreeItem[] };

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
      <div className="flex flex-row py-2 px-4 items-center justify-between mx-3">
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

      <div className="flex flex-row py-2 px-4 items-center justify-between mx-3">
        <div className="text-sm hover:text-ring">{repoFiles.author.get()}</div>
        <div className="text-sm hover:text-ring">{repoFiles.email.get()}</div>
        <div className="text-sm hover:text-ring" title={repoFiles.date.get()}>
          Last commit: {timeAgo.format(moment(repoFiles.date.get()).toDate(), 'mini')}
        </div>
        <div className="text-sm hover:text-ring">Last hash: {repoFiles.hash.get()}</div>
      </div>
      <div className="grid grid-cols-2">
        <CardContent className="px-4">
          <div className="flex flex-col gap-1">
            {repoFiles.files
              .get()
              .toSorted((a: any, b: any) => b.items?.length - a.items?.length)
              .map((item) => renderItem(item as FileTreeItem))}
          </div>
          <div className="flex flex-col gap-1"></div>
        </CardContent>
      </div>
    </div>
  );
}

const renderItem = (fileItem: FileTreeItem) => {
  if ('items' in fileItem && fileItem.items.length > 0) {
    return (
      <Collapsible key={fileItem.name}>
        <CollapsibleTrigger asChild>
          <Button
            variant="ghost"
            size="sm"
            className="group hover:bg-accent hover:text-accent-foreground w-full justify-start transition-none"
          >
            <ChevronRightIcon className="transition-transform group-data-[state=open]:rotate-90" />
            <FolderIcon />
            {fileItem.name}
          </Button>
        </CollapsibleTrigger>
        <CollapsibleContent className="style-lyra:ml-4 mt-1 ml-5">
          <div className="flex flex-col gap-1">
            {fileItem.items.map((child) => renderItem(child))}
          </div>
        </CollapsibleContent>
      </Collapsible>
    );
  }
  return (
    <Button
      key={fileItem.name}
      variant="link"
      size="sm"
      className="text-foreground w-full justify-start gap-2"
    >
      <FileIcon />
      <span>{fileItem.name}</span>
    </Button>
  );
};
