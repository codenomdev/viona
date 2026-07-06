import React from 'react';

import { Skeleton } from '@/components/ui/skeleton';

const Index: React.FC = () => {
  return (
    <div className="relative flex w-full flex-wrap justify-center min-h-screen">
      <div className="max-w-sm w-full my-auto mx-auto flex flex-col gap-7">
        <div className="flex items-center gap-4">
          <Skeleton className="size-16 shrink-0 rounded-full" />
          <div className="grid gap-2 w-full">
            <Skeleton className="h-10 max-w-lg w-full" />
            <Skeleton className="h-6 w-5/6" />
          </div>
        </div>
        <div className="flex items-center gap-4">
          <Skeleton className="size-16 shrink-0 rounded-full" />
          <div className="grid gap-2 w-full">
            <Skeleton className="h-8 max-w-lg w-full" />
            <Skeleton className="h-6 w-5/6" />
          </div>
        </div>
        <div className="flex items-center gap-4">
          <Skeleton className="size-16 shrink-0 rounded-full" />
          <div className="grid gap-2 w-full">
            <Skeleton className="h-8 max-w-lg w-full" />
            <Skeleton className="h-6 w-5/6" />
          </div>
        </div>
      </div>
    </div>
  );
};

export default React.memo(Index);
