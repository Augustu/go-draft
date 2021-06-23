#### NFS-PVC

1. ##### 安装 NFS 服务器

   arch

   ```bash
   yay -S rpcbind nfs-utils
   
   sudo vi /etc/exports
   /home/augustu/data *(rw,sync,no_root_squash,no_all_squash)
   
   sudo systemctl start rpcbind
   sudo systemctl start nfs-server.service
   
   showmount -e 192.168.50.10
   
   sudo mount -t nfs 192.168.50.10:/data d
   
   ```

   

2. ##### 配置 NFS-PVC

   ```bash
   kubectl apply -f class.yaml
   kubectl apply -f rbac.yaml
   kubectl apply -f deployment.yaml
   
   kubectl apply -f pvc-nfs-dynamic.yaml
   
   cd /home/augustu/data && ls
   default-pvc-nfs-kubedata-nginx-1-pvc-fbdf849f-b3ca-4cbf-82f0-ce023c03b06f
   
   vi index.html
   <!DOCTYPE html>
   <html>
   <head>
   <style>
   </style>
   </head>
   <body>
   
   <h1>Kubernetes - Webtest 1</h1>
   <p>This page is located on a dynamic persistent volume, and run on a k8s-cluster! :)</p>
   
   </body>
   </html>
   ```

   

3. ##### NGINX 测试 PVC

   ```bash
   kubectl apply -f deploy-nginx-1-k8s.yaml
   
   kubectl get service -n default
   nginx-1-service   NodePort    10.107.175.33   <none>        80:31728/TCP   16m
   
   http://127.0.0.1:31728
   ```

   



#### ref

1. https://www.debontonline.com/2020/11/kubernetes-part-11-how-to-configure.html nfs-pvc
2. https://www.yoyoask.com/?p=2192 nodeport
3. https://blog.csdn.net/mshxuyi/article/details/115102838 nfs

