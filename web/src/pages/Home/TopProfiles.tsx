import React, { useEffect, useState } from "react";
import { formatNumber } from "../../utils"
import { IProfile } from "../../interfaces";

export function TopProfiles() {
    var [profiles, setProfiles] = useState<IProfile[]>([])
    var [load, setLoad] = useState(false)

    useEffect(() => {
        fetch(`http://localhost/api/v1/top/profiles`)
            .then(res => res.json())
            .then((res: IProfile[]) => res.sort(() => 0.5 - Math.random()))
            .then((res: IProfile[]) => setProfiles(res.slice(0, 8)))
        setLoad(true)
    }, [])

    return (
        <div className="w-full flex flex-wrap my-5 transition-all duration-1000" style={{ opacity: load ? "1" : "0" }}>
            {
                profiles.map((profile, index) => {
                    return (
                        <div className='w-1/2 flex mt-5 px-5' key={index}>
                            <img src={profile.avatar_url} alt="avatar" className="w-14 h-14 rounded-full cursor-pointer" onClick={() => window.open(profile.link)} />
                            <div className="w-full flex flex-wrap ml-4">
                                <div className="w-full">
                                    <p className="text-lg font-bold cursor-pointer" onClick={() => window.open(profile.link)}>{profile.name}</p>
                                </div>
                                <div className="flex w-full justify-between text-gray-800">
                                    <span className="">{formatNumber(profile.subscriptions)}</span>
                                    <span className="">{formatNumber(profile.subscribers)}</span>
                                    <span className="">{formatNumber(profile.likes)}</span>
                                </div>
                            </div>
                        </div>
                    )
                })
            }
        </div>
    )
}