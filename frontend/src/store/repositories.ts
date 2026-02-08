import { hookstate, useHookstate } from '@hookstate/core';
import { toast } from 'sonner';
import { call } from '@/lib/api';

export const repositoriesState = hookstate<{ repos: object[] }>({
  repos: [],
});

export async function getRepos(query: string) {
  try {
    let { items } = await call<any[]>('repos');
    if (query !== '') {
      items = items.filter((c) => {
        return String(c.path || '')
          .toLowerCase()
          .includes(query.toLowerCase());
      });
    }
    repositoriesState.repos.set(items);
  } catch (error: any) {
    toast.error('Error! Cant load repos\n' + error.message);
    console.error('Error! Cant load repos\n' + error.message);
  }
}

export function useReposState() {
  return useHookstate(repositoriesState);
}
