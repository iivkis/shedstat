import React from "react"
import { useState } from "react"

interface ITopProfile {
    name: string
    avatar: string
}

export function Home() {
    var [link, setLink] = useState("https://www.google.com")
    var [topList, setTopList] = useState<ITopProfile[]>([
        {
            name: "🍹🥤🧃САМЫЙ СОК🧃🥤🍹",
            avatar: "https://avatars.mds.yandex.net/get-shedevrum/9283310/71d6f570758b11eea23ebecdd76c3b0f/orig",
        },
        {
            name: "Decktare",
            avatar: "https://avatars.mds.yandex.net/get-yapic/27503/lHmGhDUiVkxxzC6ijeYKaLrU-1/islands-retina-50",
        },
        {
            name: "🌀🍁eklerika🍁🌀",
            avatar: "https://avatars.mds.yandex.net/get-yapic/62162/FfIayFlGfMzIadr9WM09g0Dw1xQ-1/islands-retina-50",
        },
        {
            name: "⌞PRINCE⌝",
            avatar: "https://avatars.mds.yandex.net/get-yapic/29310/3bfL0bbQOqT3YNqZSSPeFkOVuDs-1/islands-retina-50",
        },
        {
            name: "🍹🥤🧃САМЫЙ СОК🧃🥤🍹",
            avatar: "https://avatars.mds.yandex.net/get-shedevrum/9283310/71d6f570758b11eea23ebecdd76c3b0f/orig",
        },
        {
            name: "AxisCat.com",
            avatar: "https://avatars.mds.yandex.net/get-shedevrum/11079990/bc538254876d11ee8746d65421cab00e/orig",
        },
        {
            name: "MARVEL/DC - GEEK MAFIA",
            avatar: "https://avatars.mds.yandex.net/get-yapic/40841/4zW1VVtOGb5ROrI80YD8jQVkOpE-1/islands-retina-50",
        },
        {
            name: "🟡Золотое сечение🟡",
            avatar: "https://avatars.mds.yandex.net/get-yapic/64336/cIO1dF6yNqwhpaT3kJGQdm0nIFU-1/islands-retina-50",
        },
        {
            name: "Лия Ру",
            avatar: "https://avatars.mds.yandex.net/get-yapic/27232/uweUxHRGyyTMYO6FcCcaxFZFE-1/islands-retina-50",
        },
        {
            name: "💛Astr_Vi💛",
            avatar: "https://avatars.mds.yandex.net/get-shedevrum/9283310/a59ddd6978ab11eeb7ccbebcc8c95f6d/orig",
        },
    ])

    return (
        <div className="w-2/5 flex flex-wrap justify-center mt-24 mx-auto">
            <span className='w-full text-center text-2xl' style={{ fontFamily: 'HandelGothic TL' }}>Шедеврум Статистика</span>
            <form className='w-full flex justify-center mt-20' onSubmit={(e) => { e.preventDefault(); alert(link) }}>
                <input
                    value={link}
                    className="w-full text-lg border-2 border-yellow-300 p-3 rounded-xl outline-none"
                    placeholder="Ссылка на профиль https://shedevrum.ai/@prince"
                    onChange={(e) => { setLink(e.target.value) }}
                />
            </form>
            <p className="w-full mt-10 text-lg text-justify">
                Шедеврум Статистика - это незаменимый инструмент для сервиса Шедеврум, который поможет вам получать детальную информацию о росте вашего профиля в течение определенного времени. С его помощью вы сможете отслеживать ключевые метрики, которые помогут вам анализировать эффективность вашего аккаунта.
            </p>

            <hr className="w-4/5 mt-10 mb-5" />

            <div className="w-full flex flex-wrap">
                {
                    topList.map((item, index) => {
                        return (
                            <div className='w-1/2 flex mt-5 px-5' key={index}>
                                <img src={item.avatar} alt={item.name} className="w-14 h-14 rounded-full border-yellow-500" />
                                <div className="w-full flex flex-wrap ml-4">
                                    <div className="w-full">
                                        <p className="text-lg font-bold">{item.name}</p>
                                    </div>
                                    <div className="flex w-full text-sm justify-between px-5">
                                        <span className="">443,443</span>
                                        <span className="">1,244,122</span>
                                        <span className="">4,332,342</span>
                                    </div>
                                </div>
                            </div>
                        )
                    })
                }
            </div>
        </div>
    )
}