import React, { useState } from "react"
import { IProfile, IProfileMetrics } from "../../interfaces.d"
import { TopProfiles } from './TopProfiles';
import { Metrics } from './Metrics';

export function Home() {
    var [link, setLink] = useState("")
    var [profile, setProfile] = useState<IProfile>()
    var [metrics, setMetrics] = useState<IProfileMetrics[]>()

    function loadProfile(event: React.FormEvent) {
        event.preventDefault()
        fetch(`http://localhost/api/v1/profile/${link}`)
            .then(res => res.json())
            .then(res => setProfile(res))
        fetch(`http://localhost/api/v1/profile/${link}/metrics`)
            .then(res => res.json())
            .then((res: IProfileMetrics[]) => {
                res.forEach((item) => {
                    let d = new Date(item.created_at)
                    item.created_at = `${d.getMonth().toString().padStart(2, '0')}.${d.getDay().toString().padStart(2, '0')}`
                })
                setMetrics(res)
            })
    }


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

            <p className="w-full mt-10 mb-5 text-lg text-justify" style={{ display: profile ? "none" : "" }}>
                Шедеврум Статистика - это незаменимый инструмент для сервиса Шедеврум, который поможет вам получать детальную информацию о росте вашего профиля в течение определенного времени.
                С помощью сервиса вы сможете отслеживать ключевые метрики, которые помогут вам анализировать эффективность вашего аккаунта.
            </p>

            <div className="w-full flex flex-wrap mt-5" style={{ display: profile ? "" : "none"}}>
                <Metrics profile={profile} metrics={metrics} />
            </div>

            <hr className="w-full mt-5" />

            <TopProfiles></TopProfiles>
        </div>
    )
}