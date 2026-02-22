import { hookstate, useHookstate } from '@hookstate/core';
import { toast } from 'sonner';
import { call } from '@/lib/api';

export const repoFilesState = hookstate<{
  files: object[];
  hash: '';
  message: '';
  author: '';
  email: '';
  date: '';
}>({
  files: [],
  hash: '',
  message: '',
  author: '',
  email: '',
  date: '',
});

export async function getRepoFiles(
  query: string,
  userName: string | undefined,
  repoName: string | undefined,
) {
  try {
    let response = await call<any[]>(`repos/files/${userName}/${repoName}`);
    let files = response.files;
    if (query !== '') {
      files = files.filter((c) => {
        return String(c.name || '')
          .toLowerCase()
          .includes(query.toLowerCase());
      });
    }

    repoFilesState.set({
      files: files,
      hash: response.hash,
      message: response.message,
      author: response.author,
      email: response.email,
      date: response.date,
    });
  } catch (error: any) {
    toast.error('Error! Cant load files\n' + error.message);
    console.error('Error! Cant load files\n' + error.message);
  }
}

export function useRepoFilesState() {
  return useHookstate(repoFilesState);
}
