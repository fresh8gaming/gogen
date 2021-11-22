// Copyright 2020-2021 Buf Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by protoc-gen-go-api. DO NOT EDIT.

package registryv1alpha1api

import (
	context "context"
	v1alpha1 "github.com/bufbuild/buf/private/gen/proto/go/buf/alpha/registry/v1alpha1"
)

// RepositoryCommitService is the Repository commit service.
type RepositoryCommitService interface {
	// ListRepositoryCommitsByBranch lists the repository commits associated
	// with a repository branch on a repository, ordered by their create time.
	ListRepositoryCommitsByBranch(
		ctx context.Context,
		repositoryOwner string,
		repositoryName string,
		repositoryBranchName string,
		pageSize uint32,
		pageToken string,
		reverse bool,
	) (repositoryCommits []*v1alpha1.RepositoryCommit, nextPageToken string, err error)
	// ListRepositoryCommitsByReference returns repository commits up-to and including
	// the provided reference.
	ListRepositoryCommitsByReference(
		ctx context.Context,
		repositoryOwner string,
		repositoryName string,
		reference string,
		pageSize uint32,
		pageToken string,
		reverse bool,
	) (repositoryCommits []*v1alpha1.RepositoryCommit, nextPageToken string, err error)
	// GetRepositoryCommitByReference returns the repository commit matching
	// the provided reference, if it exists.
	GetRepositoryCommitByReference(
		ctx context.Context,
		repositoryOwner string,
		repositoryName string,
		reference string,
	) (repositoryCommit *v1alpha1.RepositoryCommit, err error)
	// GetRepositoryCommitBySequenceID returns the repository commit matching
	// the provided sequence ID and branch, if it exists.
	GetRepositoryCommitBySequenceID(
		ctx context.Context,
		repositoryOwner string,
		repositoryName string,
		repositoryBranchName string,
		commitSequenceId int64,
	) (repositoryCommit *v1alpha1.RepositoryCommit, err error)
}
