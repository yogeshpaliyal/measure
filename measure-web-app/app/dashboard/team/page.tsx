"use client"

import Dropdown from "@/app/components/dropdown";
import { useState } from "react";

const teamMembers = [
  {
    id: 'asldkfjlk34343',
    email: 'anup@measure.sh'
  },
  {
    id: 'sldfkjsklf898',
    email: 'gandharva@measure.sh'
  },
  {
    id: 'asafdasfd9900',
    email: 'debjeet@measure.sh'
  },
  {
    id: 'bnflkjfg8989',
    email: 'abhay@measure.sh'
  },
  {
    id: 'cbcmvncmvn89898',
    email: 'vinu@measure.sh'
  },
  {
    id: 'sldkjkjdf8989',
    email: 'adwin@measure.sh'
  },
  {
    id: 'sbxcbvcv898',
    email: 'aakash@measure.sh'
  },
  {
    id: 'asdfsdgsdg90909',
    email: 'tarun@measure.sh'
  },
  {
    id: 'ckvjdfsfjh78aswe',
    email: 'abhinav@measure.sh'
  }
];

export default function Team() {
  const defaultTeamName = "Anup's team"
  const [saveTeamNameButtonDisabled, setSaveTeamNameButtonDisabled] = useState(true);

  return (
    <div className="flex flex-col selection:bg-yellow-200/75 items-start p-24 pt-8">
      <div className="py-4"/>
      <p className="font-display font-regular text-black text-4xl max-w-6xl text-center">Team</p>
      <div className="py-4"/>
      <p className="font-sans text-black max-w-6xl text-center">Team name</p>
      <div className="py-1"/>
      <div className="flex flex-row items-center">
        <input id="change-team-name-input" type="text" defaultValue={defaultTeamName} onChange={(event) => event.target.value === defaultTeamName? setSaveTeamNameButtonDisabled(true): setSaveTeamNameButtonDisabled(false)} className="w-96 border border-black rounded-md outline-none focus-visible:outline-yellow-300 text-black py-2 px-4 font-sans placeholder:text-neutral-400"/>
        <button disabled={saveTeamNameButtonDisabled} className="m-4 outline-none flex justify-center hover:bg-yellow-200 active:bg-yellow-300 focus-visible:bg-yellow-200 border border-black disabled:border-gray-400 rounded-md font-display text-black disabled:text-gray-400 transition-colors duration-100 py-2 px-4">Save</button>
      </div>
      <div className="py-4"/>
      <p className="font-sans text-black max-w-6xl text-center">Invite team members</p>
      <div className="py-1"/>
      <div className="flex flex-row items-center">
        <input id="invite-email-input" type="text" placeholder="Enter email" className="w-96 border border-black rounded-md outline-none focus-visible:outline-yellow-300 text-black py-2 px-4 font-sans placeholder:text-neutral-400"/>
        <div className="px-2"/>
        <Dropdown items={['Owner', 'Admin','Developer', 'Viewer']}/>
        <button className="m-4 outline-none flex justify-center hover:bg-yellow-200 active:bg-yellow-300 focus-visible:bg-yellow-200 border border-black rounded-md font-display text-black transition-colors duration-100 py-2 px-4">Invite</button>
      </div>
      <div className="py-8"/>
      <p className="font-display font-regular text-black text-2xl max-w-6xl text-center">Members</p>
      <div className="py-2"/>
      <div className="table-row-group">
        {teamMembers.map(({ id, email}) => (
            <div key={id} className="table-row font-sans text-black">
              <div className="table-cell p-4 pl-0 text-lg">{email}</div>
              <div className="table-cell p-4 w-52"><Dropdown items={['Owner', 'Admin','Developer', 'Viewer']}/></div>
              <div className="table-cell p-4"><button className="m-4 outline-none flex justify-center hover:bg-yellow-200 active:bg-yellow-300 focus-visible:bg-yellow-200 border border-black rounded-md font-display text-black transition-colors duration-100 py-2 px-4">Remove</button></div>
            </div>
        ))}
      </div>
    </div>
  )
}
