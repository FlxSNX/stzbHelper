@echo off

cd web
echo ��ǰ����Ŀ¼: %cd%
echo ��ʼ���ǰ��

call npm run build

cd ../
echo ��ǰ����Ŀ¼: %cd%
echo ��ʼ������

rem ���� dist Ŀ¼�����Ŀ¼�Ѵ�������Դ���
mkdir dist 2>nul

echo ��ʼ����windows 64λ
rem ����Ϊ Windows 64 λ��ִ���ļ�
set GOOS=windows
set GOARCH=amd64
go build -tags="nomsgpack" -ldflags="-s -w" -o dist\stzbHelper-windows-amd64.exe stzbHelper

echo ������ɣ���ִ���ļ�������� dist Ŀ¼��
pause