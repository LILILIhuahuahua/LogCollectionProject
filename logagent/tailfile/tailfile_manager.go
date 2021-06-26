package tailfile

import (
	"LogCollectionProject/logagent/common"
	"fmt"
)

/**
 * @author xhli
 * @date 2021/6/24 15:03
 * @version 1.0
 * @description: tailtask管理者
 */
type tailTaskManager struct {
	tailTaskMap map[string]*tailTask            //管理全局tailTask的map
	collectEntryList []common.CollectEntry      //初始化时，初始配置
	collectEntryChan chan []common.CollectEntry //等待新配置的通道
}

var (
	//全局的tailtask管理对象
	TTaskManager tailTaskManager
)

func InitTailTaskManager(allConf []common.CollectEntry) tailTaskManager{
	TTaskManager = tailTaskManager{
		tailTaskMap: make(map[string]*tailTask,20),
		collectEntryList:allConf,
		collectEntryChan: make(chan []common.CollectEntry),
	}
	return TTaskManager
}

//TTaskManager不断监听新的配置 (从SendNewConf不断取新的config)
func(ttManager tailTaskManager) watch(){
	//for 死循环，监听ttManager.collectEntryChan中有无新的配置
	for {
		//1、取出新的配置数组
		newConfs :=<-ttManager.collectEntryChan
		fmt.Println("get newConfs from ttManager.collectEntryChan,newConfs:",newConfs)
		//2、判断是否是新的配置
		for _,newConf :=range newConfs {
			if ttManager.isConfExist(newConf) {
				//旧的配置
				continue
			}
			//新的配置
			//根据每个收集日志项，创建并启动一个tail
			Run(ttManager,newConf)
		}

		//将newConf中不存在，而tailTakMap中存在的收集项对应的tailObj关掉
		for key,ttask :=range ttManager.tailTaskMap {
			var found bool
			for _,conf :=range newConfs{
				if key == conf.Path {
					found=true
					break
				}
			}
			//需要关闭这个任务
			if found==false{
				//ttaskManager不再管理这个任务
				delete(ttManager.tailTaskMap,key)
				//关闭这个任务的gorutine （tailfile.collectMsg的gorutine）
				ttask.cancel()
			}
		}
	}
}

func(ttManager tailTaskManager) isConfExist(conf common.CollectEntry) bool{
	_,ok :=ttManager.tailTaskMap[conf.Path]
	return ok
}
