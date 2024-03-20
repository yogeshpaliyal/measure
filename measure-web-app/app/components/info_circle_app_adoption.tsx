import React from 'react';
import { MetricsApiStatus } from '../api/api_calls';

interface InfoCircleAppAdoptionProps {
  status: MetricsApiStatus,
  value: number,
  sessions: number,
  totalSessions: number,
  title: string
}

const formatter = Intl.NumberFormat('en', { notation: 'compact' });

const InfoCircleAppAdoption = ({ status, value, sessions, totalSessions, title }: InfoCircleAppAdoptionProps) => {
  return (
    <div className="flex flex-col items-center">
      <div className={`flex flex-col group relative items-center justify-center w-64 aspect-square rounded-full border border-black border-4 hover:bg-neutral-800/25 transition-colors duration-100`}>
        {status === MetricsApiStatus.Loading && <p className="font-display text-lg">Updating...</p>}
        {status === MetricsApiStatus.Error && <p className="font-display text-lg">Error</p>}
        {status === MetricsApiStatus.Success && <p className="font-sans text-xl"> {value}%</p>}
        <div className="py-1" />
        {status === MetricsApiStatus.Success && <p className="font-sans text-sm">{formatter.format(sessions)}/{formatter.format(totalSessions)} users</p>}
        <span className="pointer-events-none z-50 max-w-xl absolute font-sans text-sm text-white rounded-md p-4 bg-neutral-800 -top-28 left-0 w-max opacity-0 transition-opacity group-hover:opacity-100">
          Adoption =  (Sessions of selected app version / Sessions of all app versions) * 100
        </span>
      </div>
      <div className="py-2" />
      <p className="font-display text-lg">{title}</p>
    </div>
  );
};

export default InfoCircleAppAdoption;