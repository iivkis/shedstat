import React, { useState } from "react";
import { formatNumber } from "../../utils"
import { IProfile } from "../../interfaces";

function PopularProfiles() {
    var [profiles, setProfiles] = useState<IProfile[]>([])
    return (
        <div className="w-full flex flex-wrap my-5">
            {
                profiles.map((item, index) => {
                    return (
                        <div className='w-1/2 flex mt-5' key={index}>
                            <img src={item.avatar_url} className="w-14 h-14 rounded-full border-yellow-500" />
                            <div className="w-full flex flex-wrap ml-4">
                                <div className="w-full">
                                    <p className="text-lg font-bold">{item.name}</p>
                                </div>
                                <div className="flex w-full text-sm px-2 text-gray-800">
                                    <span className="w-1/3 px-1">{formatNumber(item.subscriptions)}</span>
                                    <span className="w-1/3 px-1">{formatNumber(item.subscribers)}</span>
                                    <span className="w-1/3 px-1">{formatNumber(item.likes)}</span>
                                </div>
                            </div>
                        </div>
                    )
                })
            }
        </div>
    )
}