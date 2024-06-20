import { getDateAndTime } from "@utils/dateFormatter";
import { QueryOptionsData, backendApiInstance, mobileApiInstance, queryOptionsToString } from ".";
import axios from "axios";
import { formatNumber } from "@utils/formatNumber";
interface PostIncognitoUser {
    clientName: string;
    deviceId: string;
}

export const getMobileUserInfo = async (data: { userGUID: string }) => {
    const res = await mobileApiInstance.get(
      `try-get-user?userId=${data.userGUID}`
    );
    return res.data;
};
  

export const createIncognitoUser = async (data: PostIncognitoUser) => {
    const res = await mobileApiInstance.post(
      "registration/create-incognito-user/",
      data
    );
    return res.data;
};

export interface UserMobileData {
    userId: string;
    email?: string;
    name: string;
    discount: number;
    balance: number
}

export const getUserByQR = async (qr: string) => {
    const res = await mobileApiInstance.get(
        `/user-qr/try-get-user?qrContent=${qr}`
    );
    const userData: UserMobileData = {
        userId: res.data.User.UserId,
        name: res.data.User.Name,
        email: res.data.User.Email,
        discount: res.data.User.Discount*0.01,
        balance: res.data.User.ZebraCoinBalance
    }
    return userData;
}