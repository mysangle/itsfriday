import { createClient } from '@supabase/supabase-js'

let supabaseUrl;
let supabaseAnonKey;
if (typeof process !== "undefined") {
    supabaseUrl = process.env.SUPABASE_URL;
    supabaseAnonKey = process.env.SUPABASE_ANON_KEY;
} else {
    supabaseUrl = import.meta.env.VITE_SUPABASE_URL;
    supabaseAnonKey = import.meta.env.VITE_SUPABASE_ANON_KEY;
}

export const supabaseClient = createClient(supabaseUrl, supabaseAnonKey);
