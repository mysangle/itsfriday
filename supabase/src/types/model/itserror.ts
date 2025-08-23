import type { PostgrestError } from "@supabase/supabase-js";

export default class ItsError extends Error {
  code: string

  constructor(context: { message: string; code: string }) {
    super(context.message)
    this.code = context.code
  }
}

function toItsError(error: PostgrestError): ItsError {
  return new ItsError({
    message: error.message,
    code: error.code,
  })
}

export { toItsError }
