export interface User {
    id: number
    name: string
}

export interface ChatMessage {
    fromUser: string
    toUser: string
    content: string
    createdAt: string
    conversationId: string
}
