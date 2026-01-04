import { Box, ChevronRight, FolderGit, FolderGit2, Settings, ChevronDown } from 'lucide-react';
import { useLocation } from 'react-router';
import { NavLink } from 'react-router-dom';
import { useEffect, useState, useCallback } from 'react';
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from '@/components/ui/collapsible';
import {
  Sidebar,
  SidebarFooter,
  SidebarContent,
  SidebarGroup,
  SidebarMenuSub,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarMenuSubItem,
  useSidebar,
} from '@/components/ui/sidebar';
import { toast } from 'sonner';
import { cn } from '@/util';
import { useReposState, getRepos } from '@/store/repositories';

export const items = [
  {
    title: 'Main',
    url: '/',
    icon: Box,
    submenu: [],
  },
  {
    title: 'Repositories',
    icon: FolderGit,
    url: '',
    submenu: [
      // { title: 'Nodes', icon: PcCase, url: '/resource/Node' },
      // { title: 'Events', icon: CalendarSync, url: '/resource/Event' },
      // { title: 'Namespaces', icon: SquareAsterisk, url: '/resource/Namespace' },
    ],
  },
  {
    title: 'Settings',
    icon: Settings,
    url: '/settings',
    submenu: [],
  },
];

export function AppSidebar() {
  const { state } = useSidebar();
  const [sidebarItems, setSidebarItems] = useState<any>([]);
  const repos = useReposState();

  useEffect(() => {
    const reposMap = new Map<string, any>();
    repos.get().repos.forEach((r: any) => {
      reposMap.set(r.path, r);
    });
    const submenu = [...reposMap.keys()].map((k) => {
      return {
        title: k.replace(/\.git/, ''),
        icon: FolderGit2,
        url: `/resource/repositories/${reposMap.get(k).user.name}/${k}`,
        submenu: [],
      };
    });
    setSidebarItems(
      items.map((x) => {
        if (x.title === 'Repositories') {
          return {
            title: x.title,
            icon: x.icon,
            url: x.url,
            submenu: [...submenu],
          };
        } else {
          return x;
        }
      }),
    );
  }, [repos]);
  return (
    <Sidebar collapsible="offcanvas">
      <SidebarContent>
        <SidebarGroup>
          <SidebarMenu>
            {sidebarItems.map((item: any) => (
              <Collapsible key={item.title} className="group/collapsible">
                <SidebarMenuItem
                  className={cn('text-xs', state === 'collapsed' ? 'hidden' : '')}
                  key={item.title}
                >
                  {item?.url ? (
                    <MenuItem item={item} />
                  ) : (
                    <CollapsibleTrigger className="w-full" asChild>
                      <SidebarMenuButton>
                        <div className="flex flex-row items-center">
                          <item.icon size={18} className="mr-1" />
                          <div className="text-xs">{item.title}</div>
                        </div>
                        <ChevronDown
                          size={20}
                          className="ml-auto transition-transform group-data-[state=open]/collapsible:rotate-180"
                        />
                      </SidebarMenuButton>
                    </CollapsibleTrigger>
                  )}
                  <CollapsibleContent>
                    <SidebarMenuSub className="gap-0 mx-0 px-0 border-none">
                      {item.submenu.map((i, index) => (
                        <SubmenuItem key={index} item={i} />
                      ))}
                    </SidebarMenuSub>
                  </CollapsibleContent>
                </SidebarMenuItem>
              </Collapsible>
            ))}
          </SidebarMenu>
        </SidebarGroup>
      </SidebarContent>
      <SidebarFooter>
        <SidebarMenu>
          <SidebarMenuItem>
            <SidebarMenuButton>
              <div className="flex flex-row items-center ms:invisible">
                <FolderGit size={18} className="mr-1" />
                <a
                  target="_blank"
                  href="https://github.com/roman-kiselenko/smolgit"
                  className={cn('text-xs', state === 'collapsed' ? 'hidden' : '')}
                >
                  <span>smolgit v0.1.2</span>
                </a>
              </div>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarFooter>
    </Sidebar>
  );
}

function SubmenuItem({ item }: { item: any }) {
  let location = useLocation();
  if (item?.submenu?.length) {
    return (
      <SidebarMenuItem>
        <Collapsible className="group/collapsible [&[data-state=open]>button>svg:first-child]:rotate-90">
          <CollapsibleTrigger asChild>
            <SidebarMenuButton>
              <ChevronRight size={20} className="transition-transform" />
              <span className="text-xs">{item.title}</span>
            </SidebarMenuButton>
          </CollapsibleTrigger>
          <CollapsibleContent>
            <SidebarMenuSub>
              {item.submenu.map((subItem, index) => (
                <MenuItem key={index} item={subItem} />
              ))}
            </SidebarMenuSub>
          </CollapsibleContent>
        </Collapsible>
      </SidebarMenuItem>
    );
  }
  return (
    <SidebarMenuSubItem className="flex w-full" key={item.title}>
      <SidebarMenuButton className="p-1" isActive={location.pathname === item?.url}>
        <NavLink
          to={item.url}
          className="peer/menu-button flex flex-row w-full items-center gap-0 overflow-hidden rounded-md text-left outline-hidden ring-sidebar-ring transition-[width,height,padding] focus-visible:ring-2 active:bg-sidebar-accent active:text-sidebar-accent-foreground disabled:pointer-events-none disabled:opacity-50 group-has-data-[sidebar=menu-action]/menu-item:pr-8 aria-disabled:pointer-events-none aria-disabled:opacity-50 data-[active=true]:bg-sidebar-accent data-[active=true]:font-medium data-[active=true]:text-sidebar-accent-foreground data-[state=open]:hover:bg-sidebar-accent data-[state=open]:hover:text-sidebar-accent-foreground group-data-[collapsible=icon]:size-8! group-data-[collapsible=icon]:p-2! [&>span:last-child]:truncate [&>svg]:size-4 [&>svg]:shrink-0 hover:bg-sidebar-accent hover:text-sidebar-accent-foreground h-8 text-xs !gap-0"
        >
          <item.icon size={18} className="mr-1 text-gray-500 ml-2" />
          <div className="text-xs">{item.title}</div>
        </NavLink>
      </SidebarMenuButton>
    </SidebarMenuSubItem>
  );
}

function MenuItem({ item }: { item: any }) {
  let location = useLocation();
  return (
    <SidebarMenuButton isActive={location.pathname === item?.url}>
      <NavLink to={item.url} className="flex flex-row w-full items-center">
        {item?.icon ? <item.icon size={18} className="mr-1" /> : <></>}
        <div className="text-xs">{item.title}</div>
      </NavLink>
    </SidebarMenuButton>
  );
}
