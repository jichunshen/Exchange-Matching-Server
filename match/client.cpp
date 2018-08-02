#include <iostream>
#include <fstream>
#include <cstdlib>
#include <string>
#include <cstring>
#include <sys/socket.h>
#include <sys/types.h>
#include <sys/stat.h>
#include <unistd.h>
#include <netinet/in.h>
#include <arpa/inet.h>
#include <sys/socket.h>
#include <netdb.h>


int main(int argc, char ** argv){
	char* s_hostname;
	s_hostname = argv[1];
	int port_num = atoi(argv[2]);
	struct hostent * host_info = gethostbyname(s_hostname);
  	if ( host_info == NULL ) {
    	exit(EXIT_FAILURE);
  	}
  	int client_fd = socket(AF_INET,SOCK_STREAM,0);
  	struct sockaddr_in server_in;
  	server_in.sin_family = AF_INET;
  	server_in.sin_port = htons(port_num);
	memcpy(&server_in.sin_addr, host_info->h_addr_list[0], host_info->h_length);
	
	int connect_status = connect(client_fd,(struct sockaddr *)&server_in,sizeof(server_in));
	if(connect_status==-1){
		std::cout<<"ERROR: FAIL TO CONNECT TO THE SERVER\r\n";
	}
	std::cout<<"connect success\r\n";


	char buffer[256];
	memset(buffer,0,sizeof(buffer));
	std::ifstream in;
	in.open(argv[3],std::ifstream::in);
	if(!in.is_open()){
		std::cout<<"error open file\r\n";
		exit(1);
	}
	while(!in.eof()){
		memset(buffer,0,sizeof(buffer));
		in.getline(buffer,256);
		std::string temp(buffer);
		buffer[temp.length()]='\n';
		int label = send(client_fd,buffer,temp.length()+1,0);
		if(label>=0){
			std::cout<<"sent "<<label<<" bytes "<<buffer<<"\r\n";
		}
		else{
			std::cout<<"send error\r\n";
		}	
	}
	in.close();

	

	char buffer2[1024];
	memset(buffer2,0,sizeof(buffer2));
	int label = recv(client_fd,&buffer2,sizeof(buffer2),0);
	if(label>0){
		std::cout<<"recv "<<buffer2<<"\r\n";
	}
	else{
		std::cout<<"recv error\r\n";
	}	

	close(client_fd);
	return EXIT_SUCCESS;
}