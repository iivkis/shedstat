import React, { useState } from "react"
import { IProfile, IMetricsChart } from "../../interfaces.d"
import { CartesianGrid, Line, LineChart, XAxis, YAxis, Tooltip } from "recharts"

import {
    ValueType,
    NameType,
} from 'recharts/types/component/DefaultTooltipContent';

import { TooltipProps } from 'recharts';

interface ITopProfile {
    name: string
    avatar: string
    subscriptions: number
    subscribers: number
    likes: number
}

export function Home() {
    var [link, setLink] = useState("")
    var [topList, setTopList] = useState<ITopProfile[]>([
        {
            name: "🍹🥤🧃САМЫЙ СОК🧃🥤🍹",
            avatar: "https://avatars.mds.yandex.net/get-shedevrum/9283310/71d6f570758b11eea23ebecdd76c3b0f/orig",
            subscriptions: 100_000,
            subscribers: 232_423,
            likes: 100212
        },
        {
            name: "Decktare",
            avatar: "https://avatars.mds.yandex.net/get-yapic/27503/lHmGhDUiVkxxzC6ijeYKaLrU-1/islands-retina-50",
            subscriptions: 101420,
            subscribers: 103430,
            likes: 10034
        },
        {
            name: "🌀🍁eklerika🍁🌀",
            avatar: "https://avatars.mds.yandex.net/get-yapic/62162/FfIayFlGfMzIadr9WM09g0Dw1xQ-1/islands-retina-50",
            subscriptions: 102140,
            subscribers: 104240,
            likes: 100214
        },
        {
            name: "⌞PRINCE⌝",
            avatar: "https://avatars.mds.yandex.net/get-yapic/29310/3bfL0bbQOqT3YNqZSSPeFkOVuDs-1/islands-retina-50",
            subscriptions: 105360,
            subscribers: 10360,
            likes: 10046645
        },
        {
            name: "🍹🥤🧃САМЫЙ СОК🧃🥤🍹",
            avatar: "https://avatars.mds.yandex.net/get-shedevrum/9283310/71d6f570758b11eea23ebecdd76c3b0f/orig",
            subscriptions: 10640,
            subscribers: 100546,
            likes: 10640
        },
        {
            name: "AxisCat.com",
            avatar: "https://avatars.mds.yandex.net/get-shedevrum/11079990/bc538254876d11ee8746d65421cab00e/orig",
            subscriptions: 26100,
            subscribers: 146200,
            likes: 614600
        },
        {
            name: "MARVEL/DC - GEEK MAFIA",
            avatar: "https://avatars.mds.yandex.net/get-yapic/40841/4zW1VVtOGb5ROrI80YD8jQVkOpE-1/islands-retina-50",
            subscriptions: 564100,
            subscribers: 14562400,
            likes: 1026540
        },
        {
            name: "🟡Золотое сечение🟡",
            avatar: "https://avatars.mds.yandex.net/get-yapic/64336/cIO1dF6yNqwhpaT3kJGQdm0nIFU-1/islands-retina-50",
            subscriptions: 146400,
            subscribers: 5646100,
            likes: 6463100
        },
        {
            name: "Лия Ру",
            avatar: "https://avatars.mds.yandex.net/get-yapic/27232/uweUxHRGyyTMYO6FcCcaxFZFE-1/islands-retina-50",
            subscriptions: 346100,
            subscribers: 36100,
            likes: 6654100
        },
        {
            name: "💛Astr_Vi💛",
            avatar: "https://avatars.mds.yandex.net/get-shedevrum/9283310/a59ddd6978ab11eeb7ccbebcc8c95f6d/orig",
            subscriptions: 63100,
            subscribers: 64536100,
            likes: 63100
        },
    ])


    var [profile, setProfile] = useState<IProfile>()

    var [metricsChart, setMetricsChart] = useState<IMetricsChart[]>([
        {
            date: "01.01",
            likes: 10,
            subscribers: 1000,
            subscriptions: 2000
        },
        {
            date: "01.02",
            likes: 100,
            subscribers: 1500,
            subscriptions: 2000
        },
        {
            date: "01.03",
            likes: 220,
            subscribers: 1400,
            subscriptions: 2240
        },
        {
            date: "01.04",
            likes: 10,
            subscribers: 1000,
            subscriptions: 2000
        },
        {
            date: "01.05",
            likes: 100,
            subscribers: 1500,
            subscriptions: 2000
        },
        {
            date: "01.06",
            likes: 220,
            subscribers: 1400,
            subscriptions: 2240
        },
        {
            date: "01.07",
            likes: 10,
            subscribers: 1000,
            subscriptions: 2000
        },
        {
            date: "01.08",
            likes: 100,
            subscribers: 1500,
            subscriptions: 2000
        },
        {
            date: "01.09",
            likes: 220,
            subscribers: 1400,
            subscriptions: 2240
        },
        {
            date: "01.10",
            likes: 10,
            subscribers: 1000,
            subscriptions: 2000
        },
        {
            date: "01.11",
            likes: 100,
            subscribers: 1500,
            subscriptions: 2000
        },
        {
            date: "01.12",
            likes: 220,
            subscribers: 1400,
            subscriptions: 2240
        },
        {
            date: "01.01",
            likes: 10,
            subscribers: 1000,
            subscriptions: 2000
        },
        {
            date: "01.02",
            likes: 100,
            subscribers: 1500,
            subscriptions: 2000
        },
        {
            date: "01.03",
            likes: 220,
            subscribers: 1400,
            subscriptions: 2240
        },
        {
            date: "01.04",
            likes: 10,
            subscribers: 1000,
            subscriptions: 2000
        },
        {
            date: "01.05",
            likes: 100,
            subscribers: 1500,
            subscriptions: 2000
        },
        {
            date: "01.06",
            likes: 220,
            subscribers: 1400,
            subscriptions: 2240
        },
        {
            date: "01.07",
            likes: 10,
            subscribers: 1000,
            subscriptions: 2000
        },
        {
            date: "01.08",
            likes: 100,
            subscribers: 1500,
            subscriptions: 2000
        },
        {
            date: "01.09",
            likes: 220,
            subscribers: 1400,
            subscriptions: 2240
        },
        {
            date: "01.10",
            likes: 10,
            subscribers: 1000,
            subscriptions: 2000
        },
        {
            date: "01.11",
            likes: 100,
            subscribers: 1500,
            subscriptions: 2000
        },
        {
            date: "01.24",
            likes: 220,
            subscribers: 1400,
            subscriptions: 2240
        },
    ])


    function loadProfile(event: React.FormEvent) {
        event.preventDefault()
        fetch(`http://localhost:80/profile/${link}`)
            .then(res => res.json())
            .then(res => setProfile(res))
        fetch(`http://localhost:80/profile/${link}/metrics`)
            .then(res => res.json())
            .then((res: IMetricsChart[]) => {
                res.forEach((item) => {
                    let d = new Date(item.date)
                    item.date = `${d.getMonth().toString().padStart(2, '0')}.${d.getDay().toString().padStart(2, '0')}`
                })
                setMetricsChart(res)
            })
    }

    function formatNumber(num: number | ValueType | undefined): string {
        if (typeof num === 'number') {
            return num.toLocaleString('en-US')
        }
        return ""
    }

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
                    <p className="font-bold">{payload[0].payload.date}</p>
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
        <div className="w-2/5 flex flex-wrap justify-center mt-24 mx-auto">
            <span className='w-full text-center text-2xl' style={{ fontFamily: 'HandelGothic TL' }}>Шедеврум Статистика</span>
            <form className='w-full flex justify-center mt-20' onSubmit={loadProfile}>
                <input
                    value={link}
                    className="w-full text-lg border-2 border-yellow-300 p-3 rounded-xl outline-none"
                    placeholder="Ссылка на профиль https://shedevrum.ai/@prince"
                    onChange={(e) => { setLink(e.target.value) }}
                />
            </form>
            {
                !profile?.id &&
                <p className="w-full mt-10 text-lg text-justify">
                    Шедеврум Статистика - это незаменимый инструмент для сервиса Шедеврум, который поможет вам получать детальную информацию о росте вашего профиля в течение определенного времени. С его помощью вы сможете отслеживать ключевые метрики, которые помогут вам анализировать эффективность вашего аккаунта.
                </p>
            }
            {
                profile?.id &&
                <div className="w-full flex flex-wrap mt-5">
                    <div className="w-full flex justify-center mt-5">
                        <img src={profile.avatar_url} className="rounded-full w-24 h-24" />
                    </div>
                    <div className="w-full flex justify-center mt-5">
                        <span className="text-3xl font-bold tracking-tight">{profile.name}</span>
                    </div>
                    <div className="w-full flex justify-center mt-5">
                        <div className="flex flex-wrap text-center">
                            <span className="w-full text-xl font-bold">
                                {formatNumber(profile.subscriptions)}
                            </span>
                            <span className="w-full font-bold text-sm text-gray-600">
                                подписки
                            </span>
                        </div>
                        <div className="flex flex-wrap text-center">
                            <span className="w-full text-xl font-bold">
                                {formatNumber(profile.subscribers)}
                            </span>
                            <span className="w-full font-bold text-sm text-gray-600">
                                подписчиков
                            </span>
                        </div>
                        <div className="flex flex-wrap text-center">
                            <span className="w-full text-xl font-bold">
                                {formatNumber(profile.likes)}
                            </span>
                            <span className="w-full font-bold text-sm text-gray-600">
                                лайки
                            </span>
                        </div>
                    </div>
                    <LineChart width={760} height={300} margin={{ top: 10, right: 10, bottom: 10, left: 10 }} data={metricsChart}>
                        <CartesianGrid stroke="#ddd" strokeDasharray="3 3" />
                        <Tooltip content={<CustomTooltip />}></Tooltip>
                        <XAxis dataKey="date" />
                        <YAxis />
                        <Line type="monotone" dataKey="subscriptions" stroke="#303f9f" />
                    </LineChart>
                    <LineChart width={760} height={300} margin={{ top: 10, right: 10, bottom: 10, left: 10 }} data={metricsChart}>
                        <CartesianGrid stroke="#ddd" strokeDasharray="3 3" />
                        <Tooltip content={<CustomTooltip />}></Tooltip>
                        <XAxis dataKey="date" />
                        <YAxis />
                        <Line type="monotone" dataKey="subscribers" stroke="#00796b" />
                    </LineChart>
                    <LineChart width={760} height={300} margin={{ top: 10, right: 10, bottom: 10, left: 10 }} data={metricsChart}>
                        <CartesianGrid stroke="#ddd" strokeDasharray="3 3" />
                        <Tooltip content={<CustomTooltip />}></Tooltip>
                        <XAxis dataKey="date" />
                        <YAxis />
                        <Line type="monotone" dataKey="likes" stroke="#d32f2f" />
                    </LineChart>
                </div>
            }

            <hr className="w-4/5 mt-5" />
{/* 
            <div className="w-full flex flex-wrap my-5">
                {
                    topList.map((item, index) => {
                        return (
                            <div className='w-1/2 flex mt-5' key={index}>
                                <img src={item.avatar} alt={item.name} className="w-14 h-14 rounded-full border-yellow-500" />
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
            </div> */}
        </div>
    )
}