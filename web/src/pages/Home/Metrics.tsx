import React, { useState } from "react"
import { IProfile, IProfileMetrics } from "../../interfaces.d"
import { CartesianGrid, Line, LineChart, XAxis, YAxis, Tooltip } from "recharts"
import {
    ValueType,
    NameType,
} from 'recharts/types/component/DefaultTooltipContent';
import { TooltipProps } from 'recharts';
import { formatNumber } from "../../utils"

interface Props {
    profile: IProfile | undefined
    metrics: IProfileMetrics[] | undefined
}

export function Metrics({ profile, metrics }: Props) {
    function getCustomTooltipLabel(label: string | number | undefined): string {
        if (label === "subscriptions")
            return "Подписки"
        if (label === "subscribers")
            return "Подписчиков"
        if (label === "likes")
            return "Лайки"
        return ""
    }

    const CustomTooltip = ({ active, payload }: TooltipProps<ValueType, NameType>) => {
        if (active && payload && payload.length) {
            return (
                <div className="bg-white border-2 border-gray-100 p-2">
                    <p className="font-bold">{payload[0].payload.created_at}</p>
                    {
                        payload.map((entry, index) => {
                            return (
                                <div key={index} className="w-full">
                                    <div className="flex">
                                        <div style={{ color: entry.color }}>
                                            {getCustomTooltipLabel(entry.dataKey)}:
                                        </div>
                                        <div className="text-gray-600 ml-1">
                                            {formatNumber(entry.value)}
                                        </div>
                                    </div>
                                </div>
                            )
                        })
                    }
                </div>
            );
        }

        return null;
    };

    return (
        <>
            <div className="w-full flex justify-center mt-5">
                <img src={profile?.avatar_url} className="rounded-full w-24 h-24 cursor-pointer" onClick={() => window.open(profile?.link)} />
            </div>
            <div className="w-full flex justify-center mt-5">
                <span className="text-3xl font-bold cursor-pointer" onClick={() => window.open(profile?.link)}>{profile?.name}</span>
            </div>
            <div className="w-full flex justify-center mt-5">
                <div className="flex flex-wrap text-center">
                    <span className="w-full text-xl font-bold">
                        {formatNumber(profile?.subscriptions)}
                    </span>
                    <span className="w-full font-bold text-sm text-gray-600">
                        подписки
                    </span>
                </div>
                <div className="flex flex-wrap text-center">
                    <span className="w-full text-xl font-bold">
                        {formatNumber(profile?.subscribers)}
                    </span>
                    <span className="w-full font-bold text-sm text-gray-600">
                        подписчиков
                    </span>
                </div>
                <div className="flex flex-wrap text-center">
                    <span className="w-full text-xl font-bold">
                        {formatNumber(profile?.likes)}
                    </span>
                    <span className="w-full font-bold text-sm text-gray-600">
                        лайки
                    </span>
                </div>
            </div>
            <div className="w-full mt-10"></div>

            <h2 className="w-full font-bold text-xl">Подписки</h2>
            <LineChart width={760} height={300} margin={{ top: 10, right: 10, bottom: 10, left: 10 }} data={metrics}>
                <CartesianGrid stroke="#ddd" strokeDasharray="3 3" />
                <Tooltip content={<CustomTooltip />}></Tooltip>
                <XAxis dataKey="created_at" />
                <YAxis />
                <Line type="monotone" dataKey="subscriptions" stroke="#303f9f" />
            </LineChart>

            <h2 className="w-full font-bold text-xl">Подписчиков</h2>
            <LineChart width={760} height={300} margin={{ top: 10, right: 10, bottom: 10, left: 10 }} data={metrics}>
                <CartesianGrid stroke="#ddd" strokeDasharray="3 3" />
                <Tooltip content={<CustomTooltip />}></Tooltip>
                <XAxis dataKey="created_at" />
                <YAxis />
                <Line type="monotone" dataKey="subscribers" stroke="#00796b" />
            </LineChart>

            <h2 className="w-full font-bold text-xl">Лайки</h2>
            <LineChart width={760} height={300} margin={{ top: 10, right: 10, bottom: 10, left: 10 }} data={metrics}>
                <CartesianGrid stroke="#ddd" strokeDasharray="3 3" />
                <Tooltip content={<CustomTooltip />}></Tooltip>
                <XAxis dataKey="created_at" />
                <YAxis />
                <Line type="monotone" dataKey="likes" stroke="#d32f2f" />
            </LineChart>
        </>
    )
}