function prepare_and_run_docker_container {
    curr_dir=$(pwd)
    if [[ $curr_dir == *"scripts"* ]]; then
        echo "Don't run this script file inside scripts directory.
Run the script from the root directory.
Example:
    chmod +x ./scripts/run_docker.sh
    ./scripts/run_docker.sh"
        return
    fi

    # create the source directory 
    mkdir target

    # change the .config/dev.yml
    dev_config_content="$(cat .config/dev.yml)"
    source_to_change="/Users/bilginyuksel/ps-workspace/remote-code-execution/target"
    desired_source="$(pwd)/target"
    echo "${dev_config_content/source_to_change/"$desired_source"}" > .config/dev.yml

    docker build --progress=plain -t rce-engine -f prod.Dockerfile .

    # run the docker container
    docker run -it \
        -p 8888:8888 \
        -v /var/run/docker.sock:/var/run/docker.sock \
        --mount type=bind,source=$(pwd)/target,target=/target \
        rce-engine
} 

prepare_and_run_docker_container