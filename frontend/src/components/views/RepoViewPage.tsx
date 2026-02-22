import { Button } from '@/components/ui/button';
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from '@/components/ui/collapsible';
import { ChevronRightIcon, FileIcon, FolderIcon } from 'lucide-react';
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

type FileTreeItem = { name: string } | { name: string; items: FileTreeItem[] };

export function RepoViewPage() {
  const { user, repoPath } = useParams<{ user: string; repoPath: string }>();
  const repos = useReposState();
  const repoFiles = useRepoFilesState();
  const [searchQuery, setSearchQuery] = useState('');
  const { logout, AuthDisabled } = useAuth();
  const fileTree: FileTreeItem[] = [
    {
      name: 'components',
      items: [
        {
          name: 'ui',
          items: [
            { name: 'button.tsx' },
            { name: 'card.tsx' },
            { name: 'dialog.tsx' },
            { name: 'input.tsx' },
            { name: 'select.tsx' },
            { name: 'table.tsx' },
          ],
        },
        { name: 'login-form.tsx' },
        { name: 'register-form.tsx' },
      ],
    },
    {
      name: 'lib',
      items: [{ name: 'utils.ts' }, { name: 'cn.ts' }, { name: 'api.ts' }],
    },
    {
      name: 'hooks',
      items: [
        { name: 'use-media-query.ts' },
        { name: 'use-debounce.ts' },
        { name: 'use-local-storage.ts' },
      ],
    },
    {
      name: 'types',
      items: [{ name: 'index.d.ts' }, { name: 'api.d.ts' }],
    },
    {
      name: 'public',
      items: [{ name: 'favicon.ico' }, { name: 'logo.svg' }, { name: 'images' }],
    },
    { name: 'app.tsx' },
    { name: 'layout.tsx' },
    { name: 'globals.css' },
    { name: 'package.json' },
    { name: 'tsconfig.json' },
    { name: 'README.md' },
    { name: '.gitignore' },
  ];

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
          <div className="flex flex-col gap-1">{fileTree.map((item) => renderItem(item))}</div>
          {/* <DataTable
            menuDisabled={true}
            kind="repo"
            noResult={false}
            columns={columns as any}
            data={repoFiles.files.get() as any}
          /> */}
        </div>
      </div>
    </div>
  );
}

const renderItem = (fileItem: FileTreeItem) => {
  if ('items' in fileItem) {
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
