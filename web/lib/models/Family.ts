import { MapType } from "./Common";

export type Family = {
    family_id: string

    // The key is the user_id of the child user
    children: MapType<FamilyMember>

    // The key is the user_id of the parent user
    parents: MapType<FamilyMember>
};

export type FamilyMember = {
    child_call_name: string
    email: string
    name: string
    user_id: string
};