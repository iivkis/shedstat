export interface IProfile {
    id: number
    shedevrum_id: string
    created_at: string
    link: string
    name: string
    avatar_url: string
    subscriptions: number
    subscribers: number
    likes: number
}

export interface IMetricsChart {
    date: string
    subscriptions: number
    subscribers: number
    likes: number
}