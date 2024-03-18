# kube-exec-all

execute script on all pods

## RBAC

```yaml
# create serviceaccount
apiVersion: v1
kind: ServiceAccount
metadata:
  name: kube-exec-all
  namespace: autoops
---
# create clusterrole
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRole
metadata:
  name: kube-exec-all
rules:
  - apiGroups: [""]
    resources: ["pods"]
    verbs: ["list"]
  - apiGroups: [""]
    resources: ["pods/exec"]
    verbs: ["create"]
---
# create clusterrolebinding
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRoleBinding
metadata:
  name: kube-exec-all
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: kube-exec-all
subjects:
  - kind: ServiceAccount
    name: kube-exec-all
    namespace: autoops
---
# create configmap
apiVersion: v1
kind: ConfigMap
metadata:
  # !!!CHANGE ME!!!
  name:  kube-exec-all-demo
  namespace: autoops
data:
  script.sh: |
    echo hello world
---
# create cronjob
apiVersion: batch/v1beta1
kind: CronJob
metadata:
  # !!!CHANGE ME!!!
  name: kube-exec-all-demo
  namespace: autoops
spec:
  schedule: "*/5 * * * *"
  jobTemplate:
    spec:
      template:
        spec:
          serviceAccount: kube-exec-all
          containers:
            - name: kube-exec-all-demo
              image: yankeguo/kube-exec-all
              volumeMounts:
                - mountPath: /data/kube-exec-all
                  name: vol-script
          restartPolicy: OnFailure
          volumes:
            - name: vol-script
              configMap:
                # !!!CHANGE ME!!!
                name: kube-exec-all-demo
```

## Credits

GUO YANKE, MIT License
