import { hookstate, useHookstate } from '@hookstate/core';
import { toast } from 'sonner';
import { call } from '@/lib/api';

export const repoCommitsState = hookstate<{
  commits: CommitItem[];
  title: '';
  repo: object;
}>({
  commits: [],
  title: '',
  repo: {},
});

export type CommitItem = { message: string; author: string; date: string; hash: string };

export async function getRepoCommits(userName: string | undefined, repoName: string | undefined) {
  try {
    let response = await call<any>(`repos/commits/${userName}/${repoName}`);

    repoCommitsState.set({
      commits: response.commits,
      title: response.title,
      repo: response.repo,
    });
  } catch (error: any) {
    toast.error('Error! Cant load files\n' + error.message);
    console.error('Error! Cant load files\n' + error.message);
  }
}

export function useRepoCommitsState() {
  return useHookstate(repoCommitsState);
}
