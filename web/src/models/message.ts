import type { TChat } from "./chat"

export type TMessage = {
  type: 'message',
  user: string,
  chat: TChat
}
