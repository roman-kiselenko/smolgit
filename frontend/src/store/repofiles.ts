import { hookstate, useHookstate } from '@hookstate/core';
import { toast } from 'sonner';
import { call } from '@/lib/api';

export const repoFilesState = hookstate<{ files: object[] }>({
  files: [],
});

export async function getRepoFiles(query: string, userName: string | undefined, repoName: string | undefined) {
  try {
    let { files } = await call<any[]>(`repos/files/${userName}/${repoName}`);
    if (query !== '') {
      files = files.filter((c) => {
        return String(c.filename || '')
          .toLowerCase()
          .includes(query.toLowerCase());
      });
    }
    repoFilesState.files.set(files);
  } catch (error: any) {
    toast.error('Error! Cant load files\n' + error.message);
    console.error('Error! Cant load files\n' + error.message);
  }
}

export function useRepoFilesState() {
  return useHookstate(repoFilesState);
}
