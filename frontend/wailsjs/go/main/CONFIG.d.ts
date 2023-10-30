// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {main} from '../models';
import {common} from '../models';
import {time} from '../models';

export function GetConfigs():Promise<Array<main.TypeConfig>>;

export function GetENVConfigs():Promise<common.YamlInfo>;

export function GetRefreshTime():Promise<time.Time>;

export function GetStartTime():Promise<time.Time>;

export function InitConfig():Promise<boolean>;

export function InitEnv():Promise<boolean>;

export function SaveENVConfigs(arg1:common.YamlInfo):Promise<void>;

export function Start(arg1:main.Config):Promise<void>;

export function TestCmdExec(arg1:main.Config):Promise<void>;
